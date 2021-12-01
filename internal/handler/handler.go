package handler

import "telegram_bot/internal/vocabulary"

type Handler struct {
	repo       vocabulary.Repository
	commandMap map[string]func(vocabulary.Vocabulary)
}

func (h *Handler) SetHandler(path string, handler func(vocabulary.Vocabulary)) {
	h.commandMap[path] = handler
}

func New(repo vocabulary.Repository) *Handler {
	commands := make(map[string]func(vocabulary.Vocabulary))

	return &Handler{
		repo:       repo,
		commandMap: commands,
	}
}

func (h *Handler) Execute(command string) error {
	if handler, ok := h.commandMap[command]; ok {
		handler(h)
	}
	return nil
}
