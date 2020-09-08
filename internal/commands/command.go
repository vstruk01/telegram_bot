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
	err := sends.SendMessage("Enter word for Delete by Example\nWord-Translate", r.Chat_id)
	if err != nil {
		<-r.Ch.Done
		return
	}
	word := strings.TrimSpace(strings.ToLower(<-r.Ch.C))
	words := strings.Split(word, "-")
	if len(words) != 2 {
		sends.SendMessage("Write by Example\nWord-Translate", r.Chat_id)
		words = strings.Split(strings.TrimSpace(strings.ToLower(<-r.Ch.C)), "-")
		if len(words) != 2 {
			sends.SendMessage("HoW SMaRT you aRe", r.Chat_id)
			<-r.Ch.Done
			return
		}
	}
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
	err := sends.SendMessage("Enter by Example\nWord-Translate", r.Chat_id)
	if err != nil {
		log.Error.Println(err.Error())
		<-r.Ch.Done
		return
	}

	words := strings.Split(strings.TrimSpace(strings.ToLower(<-r.Ch.C)), "-")
	if len(words) != 2 {
		sends.SendMessage("Write by Example\nWord-Translate", r.Chat_id)
		words = strings.Split(strings.TrimSpace(strings.ToLower(<-r.Ch.C)), "-")
		if len(words) != 2 {
			sends.SendMessage("HoW SMaRT you aRe", r.Chat_id)
			<-r.Ch.Done
			return
		}
	}
	if !db.CheckWord(botStruct.RequestDb{
		Name:      r.Name,
		Word:      words[0],
		Translate: words[1],
		Chat_id:   r.Chat_id,
		Db:        r.OpenDb}) {
		sends.SendMessage("This word was written", r.Chat_id)
		<-r.Ch.Done
		return
	}
	if !db.AddWord(botStruct.RequestDb{
		Name:      r.Name,
		Word:      words[0],
		Translate: words[1],
		Chat_id:   r.Chat_id,
		Db:        r.OpenDb}) {
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
	message := "/start - start\n"
	message += "/help - show list of command\n"
	message += "/AddWord - add word\n"
	message += "/DeleteWord - delete word\n"
	message += "/WordKnow - mark as studied\n"
	message += "/RepeatNew - repeat new words\n"
	message += "/RepeatKnow - repeat learned words\n"
	message += "/ListNew - list of new words\n"
	message += "/ListKnow - list of lerned words\n"
	sends.SendMessage(message, r.Chat_id)
}

// // func CommandWordNew(r botStruct.Request) {
// // 	log.Info.Println("Command Word New")
// // 	err := sends.SendMessage("Enter Word Please", r.Chat_id)
// // 	if err != nil {
// // 		log.Error.Println(err.Error())
// // 	}
// // }

func CommandWordKnow(r botStruct.Request) {
	log.Info.Print("Command Word Know\n\n")
	r.Ch.Done <- true
	err := sends.SendMessage("Enter Word Please", r.Chat_id)
	if err != nil {
		<-r.Ch.Done
		return
	}

	word := strings.TrimSpace(strings.ToLower(<-r.Ch.C))
	translate := db.GetTranslate(botStruct.RequestDb{
		Name:      r.Name,
		Word:      word,
		Translate: "",
		Chat_id:   r.Chat_id,
		Db:        r.OpenDb})
	if translate == nil {
		sends.SendMessage("Sorry I did not find this word", r.Chat_id)
		<-r.Ch.Done
		return
	}
	sends.SendMessage("Enter translate of this word", r.Chat_id)
	answer := strings.TrimSpace(strings.ToLower(<-r.Ch.C))
	if strings.Contains(*translate, " "+answer+" ") {
		sends.SendMessage("Ok I beleive you", r.Chat_id)
		db.UpdateWordKnow(r.Name, word, answer, r.OpenDb)
		<-r.Ch.Done
		return
	}
	sends.SendMessage("You can not lie to me", r.Chat_id)
	<-r.Ch.Done
}

func CommandListNew(r botStruct.Request) {
	log.Info.Print("Command List New\n\n")
	r.Ch.Done <- true

	message, ok := db.GetWordsNew(r)
	if ok {
		sends.SendMessage(*message, r.Chat_id)
	}
	<-r.Ch.Done
}

func Translate(r botStruct.Request) {
	log.Info.Print("Translate\n\n")
	r.Ch.Done <- true
	translate := db.GetTranslate(botStruct.RequestDb{
		Name:      r.Name,
		Word:      r.Text,
		Translate: "",
		Chat_id:   r.Chat_id,
		Db:        r.OpenDb})
	if translate == nil {
		sends.SendMessage("I did not find this word", r.Chat_id)
		<-r.Ch.Done
		return
	}
	err := sends.SendMessage(r.Text+" -> "+*translate, r.Chat_id)
	if err != nil {
		log.Error.Println(err.Error())
	}
	<-r.Ch.Done
}

func CommandListKnow(r botStruct.Request) {
	log.Info.Println("Command List Know")
	r.Ch.Done <- true

	message, ok := db.GetWordsKnow(r)
	if ok {
		sends.SendMessage(*message, r.Chat_id)
	}
	<-r.Ch.Done
}

func CommandStart(r botStruct.Request) {
	log.Info.Println("Command Start")
	r.Ch.Done <- true
	sends.SendMessage("Hello dear, how are you ?\nDo you want to learn English ?\nSo let's go", r.Chat_id)
	<-r.Ch.Done
}
