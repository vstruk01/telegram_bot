package commands

import (
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	sends "github.com/vstruk01/telegram_bot/internal/sends"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
)

func CommandDeleteWord(r botStruct.Request) {
	err := sends.SendMessage("Enter word for delete", r.Chat_id)
	if err != nil {
		log.Error.Println(err.Error())
		r.Ch.Done <- true
		return
	}
	Word := <-r.Ch.C
	stmt, err := r.OpenDb.Prepare("DELETE FROM words WHERE word = ?")
	if err != nil {
		log.Error.Println(err.Error())
		r.Ch.Done <- true
		return
	}
	_, err = stmt.Exec(Word)
	log.CheckErr(err)
	r.Ch.Done <- true
}

func CommandAddWord(r botStruct.Request) {
	log.Info.Println("Command Add Word")
	err := sends.SendMessage("Enter by example\nWord-Translate", r.Chat_id)
	if err != nil {
		log.Error.Println(err.Error())
		r.Ch.Done <- true
		return
	}

	for_split := <-r.Ch.C
	words := strings.Split(for_split, "-")
	if len(words) != 2 {
		err := sends.SendMessage("Hmmm what wrong ?", r.Chat_id)
		if err != nil {
			log.Error.Println(err.Error())
		}
		r.Ch.Done <- true
		return
	}
	rows, err := r.OpenDb.Query("select word translate from words where word = ? and translate = ? ", words[0], words[1])
	if err != nil {
		log.Error.Println(err.Error())
		r.Ch.Done <- true
		return
	}
	if rows.Next() {
		sends.SendMessage("Sorry this word was writen", r.Chat_id)
		r.Ch.Done <- true
		return
	}
	rows.Close()
	stmt, err := r.OpenDb.Prepare("INSERT INTO words (name, word, translate, ok) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Error.Println(err.Error())
		r.Ch.Done <- true
		return
	}
	_, err = stmt.Exec(r.Name, words[0], words[1], 0)
	log.Info.Println("End Add Word")
	sends.SendMessage("Word Wrote", r.Chat_id)
	r.Ch.Done <- true
}

func CommandRepeatKnow(r botStruct.Request) {
	log.Info.Println("Command Repeat Know")
	r.Ch.Done <- true
}

func CommandRepeatNew(r botStruct.Request) {
	log.Info.Println("Command Repeat New")
	r.Ch.Done <- true
}

func CommandWordNew(r botStruct.Request) {
	log.Info.Println("Command Word New")
	err := sends.SendMessage("Enter Word Please", r.Chat_id)
	if err != nil {
		log.Error.Println(err.Error())
	}
}

func CommandWordKnow(r botStruct.Request) {
	log.Info.Println("Command Word Know")
	err := sends.SendMessage("Enter Word Please", r.Chat_id)
	var word, translate, answer string

	word = <-r.Ch.C
	rows, err := r.OpenDb.Query("select translate from words where word = ?", word)
	if err != nil {
		log.Error.Println(err.Error())
		r.Ch.Done <- true
		return
	}
	sends.SendMessage("Enter translate of this word", r.Chat_id)
	answer = <-r.Ch.C
	if !rows.Next() {
		sends.SendMessage("Sorry I did not find this word", r.Chat_id)
		r.Ch.Done <- true
		return
	}
	err = rows.Scan(&translate)
	if err != nil {
		log.Error.Println(err.Error())
		r.Ch.Done <- true
		return
	}
	if translate == answer {
		sends.SendMessage("Ok I belive you", r.Chat_id)
		err = rows.Close()
		if err != nil {
			log.Error.Println(err.Error())
		}
		_, err = r.OpenDb.Exec("update words set ok = 1 where name = ? and word = ? and translate = ?", r.Name, word, translate)
		if err != nil {
			log.Error.Println(err.Error())
		}
		r.Ch.Done <- true
		return
	}
	for rows.Next() {
		err = rows.Scan(&translate)
		if err != nil {
			log.Error.Println(err.Error())
			r.Ch.Done <- true
			return
		}
		if translate == answer {
			sends.SendMessage("Ok I belive you", r.Chat_id)
			err = rows.Close()
			if err != nil {
				log.Error.Println(err.Error())
			}
			_, err = r.OpenDb.Exec("update words set ok = 1 where name = ? and word = ? and translate = ?", r.Name, word, translate)
			if err != nil {
				log.Error.Println(err.Error())
			}
			r.Ch.Done <- true
			return
		}
	}
	err = rows.Close()
	if err != nil {
		log.Error.Println(err.Error())
	}
	sends.SendMessage("You can not lie to me", r.Chat_id)
	r.Ch.Done <- true
}

func CommandListNew(r botStruct.Request) {
	log.Info.Println("Command List New")
	m_words := make(map[string][]string)
	var word, translate, message string

	rows, err := r.OpenDb.Query("select word, translate from words where name = ? and ok = 0", r.Name)
	if log.CheckErr(err) {
		r.Ch.Done <- true
		return
	}
	for rows.Next() {
		rows.Scan(&word, &translate)
		m_words[word] = append(m_words[word], translate)
	}
	for k, vs := range m_words {
		message += k + " -> "
		for _, v := range vs {
			message += v + " "
		}
		message += "\n"
	}
	if message == "" {
		message += "empty :("
	}
	err = sends.SendMessage(message, r.Chat_id)
	if err != nil {
		log.Error.Println(err.Error())
		r.Ch.Done <- true
		return
	}
	r.Ch.Done <- true
}

func Translate(r botStruct.Request) {
	log.Info.Println("Translate")
	rows, err := r.OpenDb.Query("select word, translate from words where name = ? and word = ?", r.Name, r.Text)
	if err != nil {
		log.Error.Println(err.Error())
		r.Ch.Done <- true
		return
	}
	err = sends.SendWords(rows, r.Chat_id)
	if err != nil {
		log.Error.Println(err.Error())
		r.Ch.Done <- true
		return
	}
	r.Ch.Done <- true
	return
}

func CommandListKnow(r botStruct.Request) {
	log.Info.Println("Command List Know")
	m_words := make(map[string][]string)
	var word, translate, message string

	rows, err := r.OpenDb.Query("select word, translate from words where name = ? and ok > 0", r.Name)
	if log.CheckErr(err) {
		r.Ch.Done <- true
		return
	}
	for rows.Next() {
		rows.Scan(&word, &translate)
		m_words[word] = append(m_words[word], translate)
	}
	for k, vs := range m_words {
		message += k + " -> "
		for _, v := range vs {
			message += v + " "
		}
		message += "\n"
	}
	if message == "" {
		message += "empty :("
	}
	err = sends.SendMessage(message, r.Chat_id)
	if err != nil {
		log.Error.Println(err.Error())
		r.Ch.Done <- true
		return
	}
	r.Ch.Done <- true
}

func CommandStart(r botStruct.Request) {
	log.Info.Println("Command Start")
	err := sends.SendMessage("Hello dear, how are you ?\nDo you want to learn English ?\nSo let's go", r.Chat_id)
	log.CheckErr(err)
	r.Ch.Done <- true
}
