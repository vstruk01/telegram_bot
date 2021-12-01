package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Config struct {
	Host        string
	Port        int
	Name        string
	User        string
	Password    string
	AppName     string
	SourceFiles string
}

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(config *Config) (*Postgres, error) {
	open, err := sql.Open("postgres", config.PostgresDSN())
	if err != nil {
		return nil, err
	}

	err = open.Ping()
	if err != nil {
		return nil, err
	}

	return &Postgres{
		DB: open,
	}, nil
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s application_name=%s",
		c.Host, c.Port, c.Name, c.User, c.Password, c.AppName)
}

func (p Postgres) AddPhrase(phrase string, translate string) error {
	panic("implement me")
}

func (p Postgres) AddWord(word, translate string) error {
	panic("implement me")
}

func (p Postgres) DeleteWord(word string) error {
	panic("implement me")
}
