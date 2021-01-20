package routes

import (
	"net/http"
)

func (s *Server) handleIndex() http.HandlerFunc {
	tpl := s.mustSetupTemplate("templates/index.html")

	return func(w http.ResponseWriter, r *http.Request) {
		err := tpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
