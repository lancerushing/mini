package server

import (
	"html/template"
	"log"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	PgHost     string
	PgPort     int
	PgDb       string
	PgUsername string
	PgPassword string
}

type server struct {
	router      chi.Router
	layout      *template.Template
	db          *sqlx.DB
	loginAuth   *auth
	pwResetAuth *auth
}

func NewServer(config *Config) (*server, error) {

	s := &server{}

	s.db = connect(config)
	s.layout = template.Must(template.ParseFiles("server/templates/_layout.html"))
	s.routes()

	s.loginAuth = NewAuth("auth", "B6n6VjPbZNSw46f3yGfkCwhq", "qeZuRRCwXZqA7Z7eF9xVxbwF")
	s.pwResetAuth = NewAuth("auth-pw", "6bEAbq4camCsdbwANRT9pRut", "jZxrxNsSt6bDfRkTen62CCk5")

	return s, nil
}

func (s *server) mustSetupTemplate(fileName string) *template.Template {
	clone, err := s.layout.Clone()
	if err != nil {
		log.Fatal("Could not clone layout template: " + err.Error())
	}

	return template.Must(clone.ParseFiles(fileName))
}

func (s *server) GetHandler() chi.Router {
	return s.router
}
