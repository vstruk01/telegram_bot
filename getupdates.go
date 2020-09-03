package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

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

func CommandStart(r Request) error {
	err := SendMessage("Hello dear, how are you ?\nDo you want to learn English ?\nSo let's go", r.Chat_id)
	if err != nil {
		return err
	}
	return nil
}

func CommandRepeatKnow(r Request) error {
	// if err != nil {
	// 	return err
	// }
	return nil
}

func CommandRepeatNew(r Request) error {
	// if err != nil {
	// 	return err
	// }
	return nil
}

func CommandWordKnow(r Request) error {
	err := SendMessage("Enter Word Please", r.Chat_id)
	if err != nil {
		return err
	}
	return nil
}

func CommandWordNew(r Request) error {
	err := SendMessage("Enter Word Please", r.Chat_id)
	if err != nil {
		return err
	}
	return nil
}

func GetUpdate(url string, offset *int,
	stmt_list *map[string]sql.Rows,
	actions_list *map[string]int,
	functions map[string]func(Request) error) error {
	rest, err := GetMessage(url, offset)

	if err != nil || len(rest.Result) == 0 {
		return err
	}

	var r Request
	for _, update := range rest.Result {
		r.Text = update.Message.Text
		r.Name = update.Message.User.Username
		r.Chat_id = update.Message.Chat.Id
		log.Print("\n\033[1;34mName:\033[0m\t", r.Name,
			"\n\033[1;34mid:\033[0m\t", r.Chat_id,
			"\n\033[1;34mWrote:\033[0m\t", r.Text, "\n\n")

		err = CheckUser(r.Name, r.Chat_id)
		if err != nil {
			return err
		}

		function, ok := functions[r.Text]

		if ok {
			return function(r)
		} else {
			words := strings.Split(update.Message.Text, "-")
			if len(words) != 2 {
				err = TranslateWord(r)
				if err != nil {
					return err
				}
			} else {
				err = InsertWord(r.Name, words)
				if err != nil {
					SendMessage("Again", r.Chat_id)
					if err != nil {
						return err
					}
				} else {
					SendMessage("Ok", r.Chat_id)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func CheckUser(name string, chat_id int) error {
	var n string

	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;34mCheck User\033[0m\n")
		return err
	}
	rows, err := database.Query("select name from users WHERE name = ? and chat_id = ?", name, chat_id)
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
		err = AddUser(name, chat_id)
		if err != nil {
			return err
		} else {
			setButton(chat_id)
			return nil
		}
	}
	fmt.Print("\033[1;34mCheck User Ok\033[0m\n")
	return nil
}

func AddUser(name string, chat_id int) error {
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
	_, err = statement.Exec(name, chat_id)
	if err != nil {
		fmt.Print("\033[1;34mAdd User\033[0m\n")
		return err
	}
	return nil
}

func Know() {

}
