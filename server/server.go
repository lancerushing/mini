package server

import (
	"html/template"
	"log"
	"net/http"

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

	// @todo Make the keys configurable
	s.loginAuth = NewAuth("auth", "this-is-a-test-key-please-fix", "this-is-a-test-key-please-fix")
	s.pwResetAuth = NewAuth("auth-pw", "this-is-a-test-key-please-fix", "this-is-a-test-key-please-fix")

	s.routes()

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

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.GetHandler().ServeHTTP(w, r)
}
