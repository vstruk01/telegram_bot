package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

	if *offset <= restResponse.Result[0].Update_id {
		*offset = restResponse.Result[0].Update_id + 1
	}
	return restResponse, nil
}

func getUpdate(url string, offset *int) error {
	rest, err := getMessage(url, offset)

	if err != nil || len(rest.Result) == 0 {
		return nil;
	}
	text := rest.Result[0].Message.Text
	chat_id := rest.Result[0].Message.Chat.Id
	name := rest.Result[0].Message.User.Username
	if text == "/start" {
		sendMessage("Hello dear how are you ? Do you want learn English ? So let`s go", chat_id)
	} else if text == "/words new" {

	} else if text == "/words know" {

	} else if text == "/list" {
		respondPrintWords(name, chat_id)
	} else if text == "/make group" {

	} else if text == "/mark+" {

	} else if text == "/mark-" {

	} else {
		database, _ := sql.Open("sqlite3", "./words.db")
		for _, update := range rest.Result {
			fmt.Println("add word")
			words := strings.Split(update.Message.Text, "-")
			if len(words) != 2 {
				sendMessage("ups invalid word \n ->\t"+"\""+update.Message.Text+"\"", chat_id)
				fmt.Println("ups invalid word")
			} else {
				sendMessage("word wrote", chat_id)
				fmt.Println("word wrote")
				//CREATE TABLE IF NOT EXISTS words (id INTEGER PRIMARY KEY, name TEXT, word TEXT, transcription TEXT, translate TEXT, ok int)
				statement, _ := database.Prepare("insert into words (name, word, translate, ok)values(?, ?, ?, ?)")
				statement.Exec(name, words[0], words[1], 0)
			}
		}
	}
	return nil
}
