package manager

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vstruk01/telegram_bot/internal/commands"
	"github.com/vstruk01/telegram_bot/internal/sends"
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
	Update_id int
	Message   Message `json:"message"`
}

type RestResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

func GetMessage(url string, offset *int) (RestResponse, error) {
	resp, err := http.Get(url + "getUpdates" + "?offset=" + strconv.Itoa(*offset))

	if err != nil {
		fmt.Print("\n\033[1;34getMessage\033[0m\t")
		return RestResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print("\n\033[1;34getMessage\033[0m\t")
		return RestResponse{}, err
	}

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		fmt.Print("\n\033[1;34getMessage\033[0m\t")
		return RestResponse{}, err
	}

	if len(restResponse.Result) == 0 {
		return RestResponse{}, nil
	}

	if *offset <= restResponse.Result[0].Update_id {
		*offset = restResponse.Result[0].Update_id + 1
	}
	return restResponse, nil
}

// ! fix struct of function
func GetUpdate(master *botStruct.Master) error {
	rest, err := GetMessage(master.Url, &master.Offset)

	if err != nil || len(rest.Result) == 0 {
		return err
	}

	var r botStruct.Request
	r.OpenDb = master.OpenDb
	for _, update := range rest.Result {
		r.Text = update.Message.Text
		r.Name = update.Message.User.Username
		r.Chat_id = update.Message.Chat.Id
		log.Print()
		fmt.Print("\n\033[1;34mName:\033[0m\t\t", r.Name,
			"\n\033[1;34mChat Id:\033[0m\t", r.Chat_id,
			"\n\033[1;34mWrote:\033[0m\t\t", r.Text, "\n\n")

		err = CheckUser(r)
		if err != nil {
			return err
		}
		r.C = master.Routines[r.Chat_id]

		function, ok := master.Commands[r.Text]

		if ok {
			go function(r)
			continue
		}
		r.C <- r.Text
	}
	return nil
}

func CheckUser(r botStruct.Request) error {
	var n string

	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;34mCheck User\033[0m\n")
		return err
	}
	rows, err := database.Query("select name from users WHERE name = ? and chat_id = ?", r.Name, r.Chat_id)
	if err != nil {
		fmt.Print("\033[1;34mCheck User\033[0m\n")
		return err
	}

	if rows.Next() {
		err = rows.Scan(&n)
		err = rows.Close()
		if err != nil {
			fmt.Print("\033[1;34mCheck User\033[0m\n")
			return err
		}
	} else {
		err = AddUser(r)
		if err != nil {
			return err
		} else {
			r.C = make(chan string)
			sends.SetButton(r.Chat_id)
			return nil
		}
	}
	fmt.Print("\033[1;34mCheck User Ok\033[0m\n")
	return nil
}

func AddUser(r botStruct.Request) error {
	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;34mAdd User\033[0m\n")
		return err
	}
	statement, err := database.Prepare("insert into users (name, chat_id)values(?, ?)")
	if err != nil {
		fmt.Print("\033[1;34mAdd User\033[0m\n")
		return err
	}
	_, err = statement.Exec(r.Name, r.Chat_id)
	if err != nil {
		fmt.Print("\033[1;34mAdd User\033[0m\n")
		return err
	}
	return nil
}

func InitAll() (*botStruct.Master, error) {
	var err error
	master := new(botStruct.Master)

	// * create map of handlers
	master.Commands = make(map[string]func(botStruct.Request))
	master.HandeFunc("/start", commands.CommandStart)
	master.HandeFunc("RepeatKnow", commands.CommandRepeatKnow)
	master.HandeFunc("ListKnow", commands.CommandListKnow)
	master.HandeFunc("WordKnow", commands.CommandWordKnow)
	master.HandeFunc("WordNew", commands.CommandWordNew)
	master.HandeFunc("ListNew", commands.CommandListNew)
	master.HandeFunc("RepeatNew", commands.CommandRepeatNew)
	master.HandeFunc("DeleteWord", commands.CommandDeleteWord)

	// * initialization other veriables
	master.Url = "https://api.telegram.org/bot1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8/"
	master.Offset = 0
	master.OpenDb, err = createDB()
	if err != nil {
		return nil, err
	}

	// * create map of chans for goroutines
	master.Routines = make(map[int]chan string)
	rows, err := master.OpenDb.Query("select chat_id from users")
	if err != nil {
		fmt.Print("\033[1;32mError WordKnow = ", err.Error(), "\033[0m\n")
		return nil, err
	}
	var id int
	for rows.Next() {
		rows.Scan(&id)
		sends.SetButton(id)
		master.Routines[id] = make(chan string)
	}

	return master, nil
}

func createDB() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;34m CreateDB open\033[0m\n")
		return nil, err
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS words (id INTEGER PRIMARY KEY, name TEXT, word TEXT, transcription TEXT, translate TEXT, ok int)")
	if err != nil {
		fmt.Print("\033[1;34createDB words\033[0m\n")
		return nil, err
	}
	_, err = statement.Exec()
	if err != nil {
		fmt.Print("\033[1;34createDB open\033[0m\n")
		return nil, err
	}
	statement, err = database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, chat_id INT)")
	if err != nil {
		fmt.Print("\033[1;34mCreate DB Users\033[0m\n")
		return nil, err
	}
	_, err = statement.Exec()
	if err != nil {
		fmt.Print("\033[1;34mCreate DB Users\033[0m\n")
		return nil, err
	}
	return database, nil
}
