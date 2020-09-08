package manager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	"github.com/vstruk01/telegram_bot/internal/commands"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
	db "github.com/vstruk01/telegram_bot/internal/workdb"
)

type Chat struct {
	Id int
}

type User struct {
	Username string `json:"username"`
}

type Message struct {
	Chat Chat
	User User `json:"from"`
	Text string
}

type Update struct {
	Update_id int     `json:"update_id"`
	Message   Message `json:"message"`
}

type RestResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

func GetUpdate(master *botStruct.Master) {
	var r botStruct.Request
	r.OpenDb = master.OpenDb
	for {
		rest, err := GetMessage(&master.Offset)
		if err != nil || len(rest.Result) == 0 {
			if err != nil {
				log.Error.Println(err.Error())
			}
			continue
		}
		for _, update := range rest.Result {
			r.Text = update.Message.Text
			r.Name = update.Message.User.Username
			r.Chat_id = update.Message.Chat.Id
			log.Info.Println("\nName:\t\t", r.Name,
				"\nChat Id:\t", r.Chat_id,
				"\nWrote:\t\t", r.Text)
			if db.CheckUser(master, r.Name, r.Chat_id) != nil {
				log.Error.Println(err.Error())
				continue
			}
			r.Ch = master.Routines[r.Chat_id]
			if len(r.Ch.Done) != 0 {
				if len(r.Ch.C) == 0 {
					r.Ch.C <- r.Text
				}
				continue
			}
			function, ok := master.GetCommand(r.Text)
			if ok {
				if len(r.Ch.Done) != 0 {
					<-r.Ch.Done
				}
				go function(r)
				continue
			}
			if len(r.Ch.Done) == 0 {
				go commands.Translate(r)
				continue
			}
		}
	}
}

func GetMessage(offset *int) (RestResponse, error) {
	resp, err := http.Get(botStruct.Url + botStruct.Token + "/getUpdates" + "?offset=" + strconv.Itoa(*offset))
	if log.CheckErr(err) {
		return RestResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if log.CheckErr(err) {
		return RestResponse{}, err
	}

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if log.CheckErr(err) {
		return RestResponse{}, err
	}

	if restResponse.Ok == false {
		http.Get(botStruct.Url + botStruct.Token + "/setWebhook")
	}

	if len(restResponse.Result) == 0 {
		return RestResponse{}, nil
	}

	if *offset <= restResponse.Result[0].Update_id {
		*offset = restResponse.Result[0].Update_id + 1
	}
	return restResponse, nil
}
