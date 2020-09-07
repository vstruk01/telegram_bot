package Logger

import (
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func InitLog() {
	Info = log.New(os.Stdout, "\033[1;34mINFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "\033[1;33mWARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "\033[1;32mERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func CheckErr(err error) bool {
	if err != nil {
		Error.Println(err.Error())
		return true
	}
	return false
}
