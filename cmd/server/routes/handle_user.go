package routes

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func (s *Server) handleUserDetails() http.HandlerFunc {

	tpl := s.mustSetupTemplate("templates/user-details.html")

	return func(w http.ResponseWriter, r *http.Request) {

		userUUID := r.Context().Value(s.loginAuth.ctxKey).(string)

		user, err := s.getUser(userUUID)
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
