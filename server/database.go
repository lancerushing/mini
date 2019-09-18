package server

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func connect(config *Config) *sqlx.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		config.PgHost, config.PgPort, config.PgUsername, config.PgPassword, config.PgDb)

	db := sqlx.MustConnect("postgres", psqlInfo)

	return db
}

func (s server) getUser(uuid string) (*UserDto, error) {
	result := UserDto{}

	err := s.db.Get(&result, "SELECT uuid, name, email FROM users WHERE uuid = $1 ", uuid)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func (s server) getByEmail(email string) *UserDto {
	result := UserDto{}

	err := s.db.Get(&result, "SELECT uuid, name, email, password FROM users WHERE email = $1", email)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return &result
}
