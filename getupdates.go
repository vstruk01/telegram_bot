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
		return RestResponse{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RestResponse{}, err
	}

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
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
		return nil
	}

	var r Request
	r.Text = rest.Result[0].Message.Text
	r.Name = rest.Result[0].Message.User.Username
	r.Chat_id = rest.Result[0].Message.Chat.Id
	log.Print("\n\033[1;34mName:\033[0m\t", r.Name,
		"\n\033[1;34mid:\033[0m\t", r.Chat_id,
		"\n\033[1;34mWrote:\033[0m\t", r.Text, "\n\n\n")
	setButton(r.Chat_id)

	if CheckActions(r, stmt_list, actions_list) {
		return nil
	}
	if r.Text == "/start" {
		sendMessage("Hello dear, how are you ?\nDo you want to learn English ?\nSo let's go", r.Chat_id)
	} else if r.Text == "repeatNew" {

	} else if r.Text == "repeatKnow" {

	} else if r.Text == "listNew" {
		listNew(r.Name, r.Chat_id)
	} else if r.Text == "listKnow" {
		listKnow(r.Name, r.Chat_id)
	} else if r.Text == "Know" {
		sendMessage("Enter Word Please", r.Chat_id)
	} else if r.Text == "NotKnow" {
		sendMessage("Enter Word Please", r.Chat_id)
	} else {
		for _, update := range rest.Result {
			fmt.Println("add word")
			words := strings.Split(update.Message.Text, "-")
			if len(words) != 2 {
				sendWord(update.Message.Text, r.Name, r.Chat_id)
			} else {
				err = insertWord(r.Name, words)
				if err != nil {
					sendMessage("word did not write", r.Chat_id)
				} else {
					sendMessage("word wrote", r.Chat_id)
				}
			}
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

	} else if (*actions_list)[r.Name] == ToNotKnow {

	}
	return true
}

func Know() {

}

func insertWord(name string, words []string) error {
	var old_words string

	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		return err
	}
	rows, err := database.Query("select translate from words WHERE name = ? and word = ?", name, words[0])
	if err != nil {
		return err
	}

	rows.Next()
	err = rows.Scan(&old_words)
	rows.Close()
	if err != nil {
		statement, _ := database.Prepare("insert into words (name, word, translate, ok)values(?, ?, ?, ?)")
		statement.Exec(name, words[0], ","+words[1]+",", 0)
	} else {
		new_word := old_words + words[1] + ","
		_, err := database.Exec("update words set translate = ?1 where name = ?2 and word = ?3", new_word, name, words[0])
		if err != nil {
			return err
		}
	}
	return nil
}
