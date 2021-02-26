package manager

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	"github.com/vstruk01/telegram_bot/internal/commands"
	_ "github.com/vstruk01/telegram_bot/internal/config"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
	db "github.com/vstruk01/telegram_bot/internal/workdb"
)

type RequestDb struct {
	Name      string
	Word      string
	Translate string
	Chat_id   int
	Db        *sql.DB
}

type Request struct {
	Text    string
	Name    string
	Chat_id int
	Ch      *Channels
	OpenDb  *sql.DB
}

type Channels struct {
	C    chan string
	Done chan bool
}
type Manager struct {
	Commands map[string]func(Request) // * list command for telegram bot
	Routines map[int]*Channels        // * chanells for communication with goroutines
	Offset   int                      // * counter of request
	OpenDb   *sql.DB                  // * connect with database
}

func New() Manager {
	return Manager{}
}

func (m *Manager) HandlerFunc(command string, f func(Request)) {
	m.Commands[command] = f
}

func (m *Manager) GetUpdate() {
	var r Request
	r.OpenDb = m.OpenDb
	for {
		rest, err := Receiver(&m.Offset)
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
			if db.CheckUser(m, r.Name, r.Chat_id) != nil {
				log.Error.Println(err.Error())
				continue
			}
			r.Ch = m.Routines[r.Chat_id]
			if len(r.Ch.Done) != 0 {
				if len(r.Ch.C) == 0 {
					r.Ch.C <- r.Text
				}
				continue
			}
			function, ok := m.Commands[r.Text]
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

func Receiver(offset *int) (RestResponse, error) {
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
