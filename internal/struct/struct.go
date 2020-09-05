package botStruct

import "database/sql"

type Request struct {
	Text    string
	Name    string
	Chat_id int
	C chan string
	OpenDb  *sql.DB
}

type Master struct {
	Commands map[string]func(Request)
	Rutines  map[int]chan string
	Url      string
	Offset   int
	OpenDb   *sql.DB
}

func (m Master) HandeFunc(command string, f func(Request)) {
	m.Commands[command] = f
}

func (m Master) GetCommand(command string) (func(Request), bool) {
	f, ok := m.Commands[command]
	return f, ok
}