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

func initLog() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	Info = log.New(file, "\033[1;34mINFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(file, "\033[1;33mWARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(file, "\033[1;32mERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func CheckErr(err error) bool {
	if err != nil {
		Error.Println(err.Error())
		return true
	}
	return false
}
