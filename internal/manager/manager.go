package manager

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	"github.com/vstruk01/telegram_bot/internal/commands"
	sends "github.com/vstruk01/telegram_bot/internal/sends"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
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
			if CheckUser(r) != nil {
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
			function, ok := master.Commands[r.Text]
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

func CheckUser(r botStruct.Request) error {
	var n string

	rows, err := r.OpenDb.Query("select name from users WHERE name = ? and chat_id = ?", r.Name, r.Chat_id)
	if err != nil {
		log.Error.Println(err.Error())
		return err
	}
	if rows.Next() {
		err = rows.Scan(&n)
		err = rows.Close()
		if err != nil {
			log.Error.Println(err.Error())
			return err
		}
	} else {
		err = AddUser(r)
		if err != nil {
			log.Error.Println(err.Error())
			return err
		}
		r.Ch.C = make(chan string, 1)
		r.Ch.Done = make(chan bool, 1)
		sends.SetButton(r.Chat_id)
	}
	return nil
}

func AddUser(r botStruct.Request) error {
	statement, err := r.OpenDb.Prepare("insert into users (name, chat_id)values(?, ?)")
	if err != nil {
		log.Error.Println(err.Error())
		return err
	}
	_, err = statement.Exec(r.Name, r.Chat_id)
	if err != nil {
		log.Error.Println(err.Error())
		return err
	}
	return nil
}
