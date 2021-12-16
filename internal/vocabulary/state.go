package vocabulary

import "telegram_bot/internal/utils"

type State struct {
	successRepeated utils.Set
	needRepeat      utils.Set
	right           int64
	wrong           int64
}

func (s *State) RightAnswer(word string) {
	s.right++
	s.successRepeated[word] = struct{}{}
}

func (s *State) WrongAnswer(word string) {
	s.wrong++
	s.needRepeat[word] = struct{}{}
}
