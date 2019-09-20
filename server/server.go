package server

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// Config to hold configuration info
type Config struct {
	PgHost     string
	PgPort     int
	PgDb       string
	PgUsername string
	PgPassword string
}

// Server provides a HTTPHandler and shared resources
type Server struct {
	router      chi.Router
	layout      *template.Template
	db          *sqlx.DB
	loginAuth   *auth
	pwResetAuth *auth
	logger      *zap.Logger
}

// NewServer builds a new instance of server
func NewServer(config *Config) (*Server, error) {
	var err error
	s := &Server{}

	s.logger, err = zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer s.logger.Sync()

	s.db = connect(config)
	s.layout = template.Must(template.ParseFiles("server/templates/_layout.html"))

	// @todo Make the keys configurable
	s.loginAuth = newAuth("auth", "this-is-a-test-key-please-fix", "this-is-a-test-key-please-fix")
	s.pwResetAuth = newAuth("auth-pw", "this-is-a-test-key-please-fix", "this-is-a-test-key-please-fix")

	s.routes()

	return s, nil
}

func (s *Server) mustSetupTemplate(fileName string) *template.Template {
	clone, err := s.layout.Clone()
	if err != nil {
		log.Fatal("Could not clone layout template: " + err.Error())
	}

	return template.Must(clone.ParseFiles(fileName))
}

// GetHandler provides the Router
func (s *Server) GetHandler() http.Handler {
	return s.router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.GetHandler().ServeHTTP(w, r)
}
