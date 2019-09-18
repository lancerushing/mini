package server

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func (s *server) handleForgotPasswordForm() http.HandlerFunc {

	tpl := s.mustSetupTemplate("server/templates/forgotPasswordForm.html")

	return func(w http.ResponseWriter, r *http.Request) {

		data := map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
		}
		_ = tpl.Execute(w, data)

	}

}
