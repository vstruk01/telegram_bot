package main

type Chat struct {
	Id int
}

type Message struct {
	Chat Chat
	User User `json:"forward_from"`
	Text string
}

type User struct {
	Username string `json:"first_name"`
}

type Update struct {
	Update_id int
	Message   Message `json:"message"`
}

type RestResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type BotMessage struct {
	Chat_id int    `json:"chat_id"`
	Text    string `json:"text"`
}

type Mes struct {
	Name string
	Body string
	Time int64
}
