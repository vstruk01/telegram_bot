package commands

import (
	"strings"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	sends "github.com/vstruk01/telegram_bot/internal/sends"
	botStruct "github.com/vstruk01/telegram_bot/internal/struct"
	db "github.com/vstruk01/telegram_bot/internal/workdb"
)

func CommandDeleteWord(r botStruct.Request) {
	r.Ch.Done <- true
	log.Info.Print("Command Delete Word\n\n")
	err := sends.SendMessage("Enter word for delete by example\nWord-Translate", r.Chat_id)
	if err != nil {
		<-r.Ch.Done
		return
	}
	word := <-r.Ch.C
	words := strings.Split(word, "-")
	if !db.DeleteWord(r.Name, words[0], words[1], r.OpenDb) {
		<-r.Ch.Done
		return
	}
	sends.SendMessage("Successfully Deleted", r.Chat_id)
	log.Info.Print("Command Delete Word Ok\n\n")
	<-r.Ch.Done
}

func CommandAddWord(r botStruct.Request) {
	r.Ch.Done <- true
	log.Info.Print("Command Add Word\n\n")
	err := sends.SendMessage("Enter by example\nWord-Translate", r.Chat_id)
	if err != nil {
		log.Error.Println(err.Error())
		<-r.Ch.Done
		return
	}

	words := strings.Split(<-r.Ch.C, "-")
	if len(words) != 2 {
		sends.SendMessage("Hmmm what wrong ?", r.Chat_id)
		<-r.Ch.Done
		return
	}
	if !db.GetWord(botStruct.Request_db{r.Name, words[0], words[1], r.Chat_id, r.OpenDb}) {
		sends.SendMessage("Sorry this word was writen", r.Chat_id)
		<-r.Ch.Done
		return
	}
	if !db.AddWord(botStruct.Request_db{r.Name, words[0], words[1], r.Chat_id, r.OpenDb}) {
		<-r.Ch.Done
		return
	}
	sends.SendMessage("Word Wrote", r.Chat_id)
	log.Info.Print("Command Add Word Ok\n\n")
	<-r.Ch.Done
}

func CommandRepeatKnow(r botStruct.Request) {
	log.Info.Println("Command Repeat Know")
	r.Ch.Done <- true
	<-r.Ch.Done
}

func CommandRepeatNew(r botStruct.Request) {
	log.Info.Println("Command Repeat New")
	<-r.Ch.Done
	r.Ch.Done <- true
}

func CommandHelp(r botStruct.Request) {
	message := "/start     - початок роботи з ботом\n"
	message += "/help      - показа список команд\n"
	message += "AddWord    - додати слово або переклад\n"
	message += "DeleteWord - видалити слово\n"
	message += "WordKnow   - позначити слово як засвоєне\n"
	message += "RepeatNew  - повторити не засвоєних слова\n"
	message += "RepeatKnow - повторити засвоєні слова\n"
	message += "ListNew    - показати список не засвоєних слів\n"
	message += "ListKnow   - показати список засвоєних слів\n"
	sends.SendMessage(message, r.Chat_id)
}

// func CommandWordNew(r botStruct.Request) {
// 	log.Info.Println("Command Word New")
// 	err := sends.SendMessage("Enter Word Please", r.Chat_id)
// 	if err != nil {
// 		log.Error.Println(err.Error())
// 	}
// }

func CommandWordKnow(r botStruct.Request) {
	log.Info.Print("Command Word Know\n\n")
	r.Ch.Done <- true
	err := sends.SendMessage("Enter Word Please", r.Chat_id)
	var word, translate, answer string

	word = <-r.Ch.C
	rows, err := r.OpenDb.Query("select translate from words where word = ?", word)
	if err != nil {
		log.Error.Println(err.Error())
		<-r.Ch.Done
		return
	}
	sends.SendMessage("Enter translate of this word", r.Chat_id)
	answer = <-r.Ch.C
	if !rows.Next() {
		sends.SendMessage("Sorry I did not find this word", r.Chat_id)
		<-r.Ch.Done
		return
	}
	err = rows.Scan(&translate)
	if err != nil {
		log.Error.Println(err.Error())
		<-r.Ch.Done
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
		<-r.Ch.Done
		return
	}
	for rows.Next() {
		err = rows.Scan(&translate)
		if err != nil {
			log.Error.Println(err.Error())
			<-r.Ch.Done
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
			<-r.Ch.Done
			return
		}
	}
	err = rows.Close()
	if err != nil {
		log.Error.Println(err.Error())
	}
	sends.SendMessage("You can not lie to me", r.Chat_id)
	<-r.Ch.Done
}

func CommandListNew(r botStruct.Request) {
	log.Info.Print("Command List New\n\n")
	r.Ch.Done <- true
	m_words := make(map[string][]string)
	var word, translate, message string

	rows, err := r.OpenDb.Query("select word, translate from words where name = ? and ok = 0", r.Name)
	if log.CheckErr(err) {
		<-r.Ch.Done
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
	}
	<-r.Ch.Done
}

func Translate(r botStruct.Request) {
	log.Info.Print("Translate\n\n")
	r.Ch.Done <- true
	rows, err := r.OpenDb.Query("select word, translate from words where name = ? and word = ?", r.Name, r.Text)
	if err != nil {
		log.Error.Println(err.Error())
		<-r.Ch.Done
		return
	}
	err = sends.SendWords(rows, r.Chat_id)
	if err != nil {
		log.Error.Println(err.Error())
	}
	<-r.Ch.Done
}

func CommandListKnow(r botStruct.Request) {
	log.Info.Println("Command List Know")
	r.Ch.Done <- true
	m_words := make(map[string][]string)
	var word, translate, message string

	rows, err := r.OpenDb.Query("select word, translate from words where name = ? and ok > 0", r.Name)
	if log.CheckErr(err) {
		<-r.Ch.Done
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
	}
	<-r.Ch.Done
}

func CommandStart(r botStruct.Request) {
	log.Info.Println("Command Start")
	r.Ch.Done <- true
	err := sends.SendMessage("Hello dear, how are you ?\nDo you want to learn English ?\nSo let's go", r.Chat_id)
	log.CheckErr(err)
	<-r.Ch.Done
}
