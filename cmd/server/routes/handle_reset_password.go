package routes

import (
	"encoding/base64"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleRestPasswordForm() http.HandlerFunc {
	tpl := s.mustSetupTemplate("templates/resetPasswordForm.html")

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

// ################### Submit ###################

func (s *Server) handleRestPasswordSubmit() http.HandlerFunc {
	tplFail := s.mustSetupTemplate("templates/resetPasswordForm.html")
	tplSuccess := s.mustSetupTemplate("templates/resetPasswordSuccess.html")

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

		userUUID := r.Context().Value(s.pwResetAuth.ctxKey).(string)
		bcryptBytes, _ := bcrypt.GenerateFromPassword([]byte(pass1), bcrypt.MinCost)

		sql := "UPDATE users SET  password = :password WHERE uuid = :uuid"
		dbData := map[string]interface{}{
			"password": string(bcryptBytes),
			"uuid":     userUUID,
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

// ################### Verify incoming token ###################

func (s *Server) handleResetPasswordVerify() http.HandlerFunc {
	tplFail := s.mustSetupTemplate("templates/resetPasswordVerifyFail.html")

	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{}

		tokens, ok := r.URL.Query()["token"]

		if !ok || len(tokens[0]) < 1 {
			data["errorMsg"] = "Token is missing"

			w.WriteHeader(http.StatusBadRequest)

			err := tplFail.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		token := tokens[0]

		tokenBytes, err := base64.URLEncoding.DecodeString(token)
		if err != nil {
			data["errorMsg"] = "Token will not decode"

			w.WriteHeader(http.StatusBadRequest)

			err := tplFail.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		userUUID, err := tokenExtractMessage(tokenBytes)
		if err != nil {
			data["errorMsg"] = err.Error()

			w.WriteHeader(http.StatusBadRequest)

			err := tplFail.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		uuidString, err := uuid.FromBytes(userUUID)
		if err != nil {
			data["errorMsg"] = err.Error()

			w.WriteHeader(http.StatusBadRequest)

			err := tplFail.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		err = s.pwResetAuth.setCookie(w, uuidString.String())
		if err != nil {
			data["errorMsg"] = err.Error()

			w.WriteHeader(http.StatusBadRequest)

			err := tplFail.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}

		http.Redirect(w, r, "../", http.StatusSeeOther)
	}
}
