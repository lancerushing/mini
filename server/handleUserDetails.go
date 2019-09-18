package server

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func (s *server) handleUserDetails() http.HandlerFunc {

	tpl := s.mustSetupTemplate("server/templates/user-details.html")

	return func(w http.ResponseWriter, r *http.Request) {

		userUuid := r.Context().Value(s.loginAuth.ctxKey).(string)

		user, err := s.getUser(userUuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
			"name":           user.Name,
			"email":          user.Email,
		}

		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
