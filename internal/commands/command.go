package commands

import (
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	"github.com/vstruk01/telegram_bot/internal/sends"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
)

func CommandDeleteWord(r botStruct.Request) {
	err := sends.SendMessage("Enter please\nword for delete", r.Chat_id)
	if log.CheckErr(err) {
		r.Ch.Done <- true
		return
	}
	Word := <-r.Ch.C
	stmt, err := r.OpenDb.Prepare("DELETE FROM words WHERE word = ?")
	if log.CheckErr(err) {
		r.Ch.Done <- true
		return
	}
	_, err = stmt.Exec(Word)
	log.CheckErr(err)
	r.Ch.Done <- true
}

func CommandAddWord(r botStruct.Request) {
	err := sends.SendMessage("Enter please\nWord-Translate", r.Chat_id)
	if log.CheckErr(err) {
		r.Ch.Done <- true
		return
	}

	words := strings.Split(<-r.Ch.C, "-")
	if len(words) != 2 {
		err := sends.SendMessage("Hmmm what wrong ?", r.Chat_id)
		log.CheckErr(err)
		r.Ch.Done <- true
		return
	}
	stmt, err := r.OpenDb.Prepare("INSERT INTO words (name, word, translate, ok) VALUES(?, ?, ?, ?)")
	if log.CheckErr(err) {
		r.Ch.Done <- true
		return
	}
	_, err = stmt.Exec(r.Name, words[0], words[1], 0)
	log.CheckErr(err)
	r.Ch.Done <- true
}

func CommandRepeatKnow(r botStruct.Request) {
	r.Ch.Done <- true
}

func CommandRepeatNew(r botStruct.Request) {
	r.Ch.Done <- true
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

	word = <-r.Ch.C
	rows, err := r.OpenDb.Query("select translate from words where word = ?", word)
	if log.CheckErr(err) {
		r.Ch.Done <- true
		return
	}
	if !rows.Next() {
		sends.SendMessage("Sorry I did not find this word", r.Chat_id)
		r.Ch.Done <- true
		return
	}
	err = rows.Scan(&translate)
	if log.CheckErr(err) {
		r.Ch.Done <- true
		return
	}
	rows.Close()

	translates = strings.Split(translate, ",")
	sends.SendMessage("Enter translate of this word", r.Chat_id)
	answer = <-r.Ch.C
	for _, translate = range translates {
		if translate == answer {
			fmt.Print("\033[1;34mNice\033[0m\n")
			sends.SendMessage("Ok I belive you", r.Chat_id)
			_, err = r.OpenDb.Exec("update words set ok = 1 where name = ?1 and word = ?2", r.Name, word)
			log.CheckErr(err)
			r.Ch.Done <- true
			return
		}
	}
	sends.SendMessage("You can not lie to me", r.Chat_id)
	r.Ch.Done <- true
}

func CommandListNew(r botStruct.Request) {
	rows, err := r.OpenDb.Query("select word, translate from words where name = ? and ok = 0", r.Name)
	if log.CheckErr(err) {
		r.Ch.Done <- true
		return
	}
	sends.SendWords(rows, r.Chat_id)
	r.Ch.Done <- true
}

func CommandListKnow(r botStruct.Request) {
	rows, err := r.OpenDb.Query("select word, translate from words where name = ? and ok > 0", r.Name)
	if log.CheckErr(err) {
		r.Ch.Done <- true
		return
	}
	sends.SendWords(rows, r.Chat_id)
	r.Ch.Done <- true
}

func CommandStart(r botStruct.Request) {
	err := sends.SendMessage("Hello dear, how are you ?\nDo you want to learn English ?\nSo let's go", r.Chat_id)
	log.CheckErr(err)
	r.Ch.Done <- true
}

func Translate(r botStruct.Request) {

}
