package routes

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	// import DB into namespace.
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func connect(config *Config) *sqlx.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		config.PgHost, config.PgPort, config.PgUsername, config.PgPassword, config.PgDB)

	db := sqlx.MustConnect("postgres", psqlInfo)

	return db
}

func (s Server) getUser(uuid string) (*UserDto, error) {
	result := UserDto{}

	err := s.db.Get(&result, "SELECT uuid, name, email FROM users WHERE uuid = $1 ", uuid)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s Server) getByEmail(email string) *UserDto {
	result := []UserDto{}

	err := s.db.Select(&result, "SELECT uuid, email, password FROM users WHERE email = $1", email)
	if err != nil {
		log.Error().Err(err).Msg("Bad Query")

		return nil
	}

	if len(result) != 1 {
		return nil
	}

	return &result[0]
}

// UserDto holds info from DB.
type UserDto struct {
	UUID     string
	Email    string
	Name     string
	Password string
}
