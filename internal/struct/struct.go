package botStruct

import "database/sql"

var (
	Url   string = "https://api.telegram.org/bot"                   // * url of telegram
	Token string = "1060785017:AAG7eJUSygisjIF_g97Dj5TKVzS-ct76su8" // * your token of telegram bot
)

type Request_db struct {
	Name      string
	Word      string
	Translate string
	Chat_id   int
	Db        *sql.DB
}

type Request struct {
	Text    string
	Name    string
	Chat_id int
	Ch      *Channels
	OpenDb  *sql.DB
}

type Channels struct {
	C    chan string
	Done chan bool
}

type Master struct {
	Commands map[string]func(Request) // * list command for telegram bot
	Routines map[int]*Channels        // * chanells for communication with goroutines
	Offset   int                      // * counter of request
	OpenDb   *sql.DB                  // * connect with database
}

func (m Master) HandeFunc(command string, f func(Request)) {
	m.Commands[command] = f
}

func (m Master) GetCommand(command string) (func(Request), bool) {
	f, ok := m.Commands[command]
	return f, ok
}
