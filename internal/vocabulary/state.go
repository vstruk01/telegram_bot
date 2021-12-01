package vocabulary

import "telegram_bot/internal/utils"

type State struct {
	repeatedWords utils.Set
	right         int64
	wrong         int64
}

func (s *State) RightAnswer() {
	s.right++
}

func (s *State) WrongAnswer() {
	s.wrong++
}

func (s *State) Repeated(word string) {

}
