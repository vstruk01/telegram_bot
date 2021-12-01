package handler

type Performer interface {
	Add()
	Delete()
	Repeat()
	Execute(command string) error
}

type Repository interface {
	AddPhrase(phrase string, translate string) error
	AddWord(word, translate string) error
	DeleteWord(word string) error
}

type Handler struct {
	Repo Repository
	CommandMap map[string]func(Performer)
}


func New(repo Repository, commands map[string]func(Performer)) *Handler {
	return &Handler{
		Repo: repo,
		CommandMap: commands,
	}
}

func (h Handler) Add() {
	panic("implement me")
}

func (h Handler) Delete() {
	panic("implement me")
}

func (h Handler) Repeat() {
	panic("implement me")
}

func (h Handler) Execute(command string) error {
	if fn, ok := h.CommandMap[command]; ok {
		fn(h)
	}
	return nil
}

