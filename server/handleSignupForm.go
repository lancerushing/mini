package server

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func (s *server) handleSignupForm() http.HandlerFunc {

	tpl := s.mustSetupTemplate("server/templates/signupForm.html")

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
