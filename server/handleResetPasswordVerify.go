package server

import (
	"encoding/base64"
	"net/http"

	"github.com/google/uuid"
)

func (s *server) handleResetPasswordVerify() http.HandlerFunc {

	tplFail := s.mustSetupTemplate("server/templates/resetPasswordVerifyFail.html")

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

		userUuid, err := tokenExtractMessage(tokenBytes)
		if err != nil {
			data["errorMsg"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			err := tplFail.Execute(w, data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		uuidString, err := uuid.FromBytes(userUuid)
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
