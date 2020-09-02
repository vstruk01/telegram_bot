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

	_ "github.com/mattn/go-sqlite3"
)

func getMessage(url string, offset *int) (RestResponse, error) {
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

func getUpdate(url string, offset *int,
	stmt_list *map[string]sql.Rows,
	actions_list *map[string]int) error {
	rest, err := getMessage(url, offset)

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

		if CheckActions(r, stmt_list, actions_list) {
			return nil
		}
		if r.Text == "/start" {
			err = sendMessage("Hello dear, how are you ?\nDo you want to learn English ?\nSo let's go", r.Chat_id)
			if err != nil {
				return err
			}
		} else if r.Text == "repeatNew" {
			if err != nil {
				return err
			}
		} else if r.Text == "repeatKnow" {
			if err != nil {
				return err
			}
		} else if r.Text == "listNew" {
			err = listNew(r.Name, r.Chat_id)
			if err != nil {
				return err
			}
		} else if r.Text == "listKnow" {
			err = listKnow(r.Name, r.Chat_id)
			if err != nil {
				return err
			}
		} else if r.Text == "Know" {
			err = sendMessage("Enter Word Please", r.Chat_id)
			if err != nil {
				return err
			}
		} else if r.Text == "NotKnow" {
			err = sendMessage("Enter Word Please", r.Chat_id)
			if err != nil {
				return err
			}
		} else {
			words := strings.Split(update.Message.Text, "-")
			if len(words) != 2 {
				fmt.Print("\033[1;34mTo send Word\033[0m\t\n")
				err = sendWord(update.Message.Text, r.Name, r.Chat_id)
				if err != nil {
					return err
				}
			} else {
				err = insertWord(r.Name, words)
				if err != nil {
					sendMessage("word did not write", r.Chat_id)
					if err != nil {
						return err
					}
					fmt.Print("\033[1;34mWord did not Write\033[0m\t\n")
				} else {
					sendMessage("word wrote", r.Chat_id)
					if err != nil {
						return err
					}
					fmt.Print("\033[1;34mWord Wrote\033[0m\t\n")
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
		fmt.Print("\n\033[1;34mcheckUser\033[0m\t")
		return err
	}
	rows, err := database.Query("select name from users WHERE name = ? and chat_id = ?", name, chat_id)
	if err != nil {
		fmt.Print("\n\033[1;34mcheckUser\033[0m\t")
		return err
	}

	if rows.Next() {
		err = rows.Scan(&n)
		if err != nil {
			fmt.Print("\n\033[1;34mcheckUser\033[0m\t")
			return err
		}
		err = rows.Close()
		if err != nil {
			fmt.Print("\n\033[1;34mcheckUser\033[0m\t")
			return err
		}
		if n == "" {
			setButton(chat_id)
		}
	}
	return nil
}

func CheckActions(r Request,
	stmt_list *map[string]sql.Rows,
	actions_list *map[string]int) bool {

	if (*actions_list)[r.Name] == ToNone {
		return false
	}
	if (*actions_list)[r.Name] == ToKnow {

	} else if (*actions_list)[r.Name] == ToNotKnow {

	} else if (*actions_list)[r.Name] == ToRepeatNew {

	} else if (*actions_list)[r.Name] == ToRepeatKnow {

	}
	return true
}

func Know() {

}
