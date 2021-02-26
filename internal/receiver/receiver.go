package receiver

import (
	log "github.com/vstruk01/telegram_bot/internal/Logger"
	_ "github.com/vstruk01/telegram_bot/internal/config"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Chat struct {
	Id int
}

type User struct {
	Username string `json:"username"`
}

type Message struct {
	Chat Chat
	User User `json:"from"`
	Text string
}

type Update struct {
	Update_id int     `json:"update_id"`
	Message   Message `json:"message"`
}

type RestResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Receiver struct {
	Offset int
}

func New(offset int) Receiver {
	return Receiver{
		Offset: offset,
	}
}

func (r *Receiver) receivingMessage(offset *int) (RestResponse, error) {
	resp, err := http.Get(Url + Token + "/getUpdates" + "?offset=" + strconv.Itoa(*offset))
	if log.CheckErr(err) {
		return RestResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if log.CheckErr(err) {
		return RestResponse{}, err
	}

	restResponse := RestResponse{}
	err = json.Unmarshal(body, &restResponse)
	if log.CheckErr(err) {
		return RestResponse{}, err
	}

	if restResponse.Ok == false {
		http.Get(Url + Token + "/setWebhook")
	}

	if len(restResponse.Result) == 0 {
		return RestResponse{}, nil
	}

	if *offset <= restResponse.Result[0].Update_id {
		*offset = restResponse.Result[0].Update_id + 1
	}
	return restResponse, nil
}
