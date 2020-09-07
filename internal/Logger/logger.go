package Logger

import (
	Init "github.com/vstruk01/telegram_bot/init"
)

func CheckErr(err error) bool {
	if err != nil {
		Init.Error.Println(err.Error())
		return true
	}
	return false
}