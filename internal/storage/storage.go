package storage

type Storage interface {
	GetUsersIDs() (*[]int, error)
	CheckUser(name string, id int) error
	AddUser(name string, id int) error
	GetWordsNew(name string) (*map[string]string, error)
	GetWordsKnow(name string) (*map[string]string, error)
	GetTranslate(name, word string) (*string, error)
	CheckWord(word, translate, name string) (bool, error)
	DeleteWord(name, word, translate string) error
	AddWord(word, translate, name string) error
	UpdateWordKnow(name, word, translate string) error
	InitTables() error
}
