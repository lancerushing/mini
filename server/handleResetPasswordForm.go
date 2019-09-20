package server

import (
	"github.com/gorilla/csrf"

	"net/http"
)

func (s *Server) handleRestPasswordForm() http.HandlerFunc {

	tpl := s.mustSetupTemplate("server/templates/resetPasswordForm.html")

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
