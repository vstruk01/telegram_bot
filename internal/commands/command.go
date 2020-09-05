package commands

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vstruk01/telegram_bot/internal/sends"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
)

func CommandDeleteWord(r botStruct.Request) {

}

func CommandRepeatKnow(r botStruct.Request) {
}

func CommandRepeatNew(r botStruct.Request) {
}

func CommandWordNew(r botStruct.Request) {
	err := sends.SendMessage("Enter Word Please", r.Chat_id)
	if err != nil {
		fmt.Print("\033[1;32mError WordNew = ", err.Error(), "\033[0m\n")
	}
}

func CommandWordKnow(r botStruct.Request) {
	err := sends.SendMessage("Enter Word Please", r.Chat_id)
	var word, translate, answer string
	var translates []string

	fmt.Print("\033[1;34mWait Word\033[0m\n")
	word = <-r.C
	fmt.Print("\033[1;34mGet Word Yes\033[0m\n")
	rows, err := r.OpenDb.Query("select translate from words where word = ?", word)
	if err != nil {
		fmt.Print("\033[1;32mError WordKnow = ", err.Error(), "\033[0m\n")
		return
	}
	if !rows.Next() {
		sends.SendMessage("Sorry I do not find this word", r.Chat_id)
		return
	}
	err = rows.Scan(&translate)
	if err != nil {
		fmt.Print("\033[1;32mError WordKnow = ", err.Error(), "\033[0m\n")
		return
	}
	rows.Close()

	translates = strings.Split(translate, ",")
	sends.SendMessage("Enter translate of this word", r.Chat_id)
	fmt.Print("\033[1;34mWait Answer\033[0m\n")
	answer = <-r.C
	fmt.Print("\033[1;34mGet Answer Yes\033[0m\n")
	for _, translate = range translates {
		if translate == answer {
			fmt.Print("\033[1;34mNice\033[0m\n")
			sends.SendMessage("Ok I belive you", r.Chat_id)
			_, err = r.OpenDb.Exec("update words set ok = 1 where name = ?1 and word = ?2", r.Name, word)
			if err != nil {
				fmt.Print("\033[1;32mError WordKnow = ", err.Error(), "\033[0m\n")
			}
			return
		}
	}
	fmt.Print("\033[1;34mNot Nice\033[0m\n")
	sends.SendMessage("You can not lie to me", r.Chat_id)
}

func CommandListNew(r botStruct.Request) {
	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;32mError ListNew = ", err.Error(), "\033[0m\n")
		return
	}
	rows, err := database.Query("select word, translate from words where name = ? and ok = 0", r.Name)
	if err != nil {
		fmt.Print("\033[1;32mError ListNew = ", err.Error(), "\033[0m\n")
		return
	}
	sends.SendWords(rows, r.Chat_id)
	fmt.Print("\033[1;34mlistNew Ok\033[0m\n")
}

func CommandListKnow(r botStruct.Request) {
	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;32mError listKnow = ", err.Error(), "\033[0m\n")
		return
	}
	rows, err := database.Query("select word, translate from words where name = ? and ok > 0", r.Name)
	if err != nil {
		fmt.Print("\033[1;32mError listKnow = ", err.Error(), "\033[0m\n")
		return
	}
	sends.SendWords(rows, r.Chat_id)
	fmt.Print("\033[1;34mlistKnow Ok\033[0m\n")
}

func InsertWord(name string, words []string) error {
	var old_words string

	database, err := sql.Open("sqlite3", "./words.db")
	if err != nil {
		fmt.Print("\033[1;34mInsert Word\033[0m\n")
		return err
	}
	rows, err := database.Query("select translate from words WHERE name = ? and word = ?", name, words[0])
	defer rows.Close()
	if err != nil {
		fmt.Print("\033[1;34mInsert Word\033[0m\n")
		return err
	}

	rows.Next()
	err = rows.Scan(&old_words)
	if err != nil {
		statement, _ := database.Prepare("insert into words (name, word, translate, ok)values(?, ?, ?, ?)")
		statement.Exec(name, words[0], ","+words[1]+",", 0)
	} else {
		new_word := old_words + words[1] + ","
		_, err := database.Exec("update words set translate = ?1 where name = ?2 and word = ?3", new_word, name, words[0])
		if err != nil {
			fmt.Print("\033[1;34mInsert Word\033[0m\n")
			return err
		}
	}
	fmt.Print("\033[1;34mInsert Word Ok\033[0m\n")
	return nil
}

func CommandStart(r botStruct.Request) {
	err := sends.SendMessage("Hello dear, how are you ?\nDo you want to learn English ?\nSo let's go", r.Chat_id)
	if err != nil {
		fmt.Print("\033[1;32mError Command Start = ", err.Error(), "\033[0m\n")
	}
}
