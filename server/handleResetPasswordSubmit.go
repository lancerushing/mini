package server

import (
	"net/http"

	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) handleRestPasswordSubmit() http.HandlerFunc {

	tplFail := s.mustSetupTemplate("server/templates/resetPasswordForm.html")
	tplSuccess := s.mustSetupTemplate("server/templates/resetPasswordSuccess.html")

	return func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse: "+err.Error(), http.StatusInternalServerError)
			return
		}
		pass1 := r.FormValue("password")
		pass2 := r.FormValue("password2")

		if pass1 != pass2 || pass1 == "" {
			data := map[string]interface{}{
				csrf.TemplateTag: csrf.TemplateField(r),
				"errorMsg":       "Passwords do not match.",
			}
			w.WriteHeader(http.StatusBadRequest)
			err := tplFail.Execute(w, data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		userUuid := r.Context().Value(s.pwResetAuth.ctxKey).(string)
		bcryptBytes, _ := bcrypt.GenerateFromPassword([]byte(pass1), bcrypt.MinCost)

		sql := "UPDATE users SET  password = :password WHERE uuid = :uuid"
		dbData := map[string]interface{}{
			"password": string(bcryptBytes),
			"uuid":     userUuid,
		}

		_, err = s.db.NamedExec(sql, dbData)
		if err != nil {
			data := map[string]interface{}{
				csrf.TemplateTag: csrf.TemplateField(r),
				"errorMsg":       err.Error(),
			}
			w.WriteHeader(http.StatusBadRequest)
			err = tplFail.Execute(w, data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		s.pwResetAuth.deleteCooke(w)
		w.WriteHeader(http.StatusAccepted)
		err = tplSuccess.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}

}
