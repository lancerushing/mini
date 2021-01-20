package routes

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// Config to hold configuration info.
type Config struct {
	PgHost     string
	PgPort     int
	PgDB       string
	PgUsername string
	PgPassword string
}

// Server provides a HTTPHandler and shared resources.
type Server struct {
	router      chi.Router
	layout      *template.Template
	db          *sqlx.DB
	loginAuth   *auth
	pwResetAuth *auth
}

// NewServer builds a new instance of routes.
func NewServer(config *Config) (*Server, error) {
	var err error

	s := &Server{}

	s.db = connect(config)

	//go:embed templates/_layout.html
	var templateString string

	s.layout, err = template.New("_layout").Parse(templateString)
	if err != nil {
		return nil, err
	}

	// @todo Make the keys configurable
	s.loginAuth = newAuth("auth", "this-is-a-test-key-please-fix", "this-is-a-test-key-please-fix")
	s.pwResetAuth = newAuth("auth-pw", "this-is-a-test-key-please-fix", "this-is-a-test-key-please-fix")

	s.routes()

	return s, nil
}

//go:embed templates/*
var tmplFS embed.FS

func getContents(path string) []byte {
	data, err := tmplFS.ReadFile(path)
	if err != nil {
		log.Fatal().Err(err).Msg("could not find embedded file")
	}

	return data
}

func (s *Server) mustSetupTemplate(path string) *template.Template {
	clone, err := s.layout.Clone()
	if err != nil {
		log.Fatal().Msg("could not clone layout template: " + err.Error())
	}

	t, err := clone.Parse(string(getContents(path)))
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse content")
	}

	return t
}

// GetHandler provides the Router.
func (s *Server) GetHandler() http.Handler {
	return s.router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.GetHandler().ServeHTTP(w, r)
}
