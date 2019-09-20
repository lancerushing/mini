package server

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func (s *Server) handleLoginForm() http.HandlerFunc {

	tpl := s.mustSetupTemplate("server/templates/loginForm.html")

	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
		}

		err := tpl.Execute(w, data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}
