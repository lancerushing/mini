package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CustomPayload struct {
	jwt.Payload
	Uuid string `json:"Uuid,omitempty"`
}

func (s *server) handleForgotPasswordSubmit() http.HandlerFunc {
	tplSuccess := s.mustSetupTemplate("server/templates/forgotPasswordSuccess.html")

	tplEmailHtml := template.Must(template.ParseFiles("server/templates/forgotPasswordEmail.html"))
	tplEmailText := template.Must(template.ParseFiles("server/templates/forgotPasswordEmail.text"))

	sendResetLink := func(email string) error {
		if len(email) == 0 {
			return errors.Errorf("Email is empty")
		}

		existingUser := s.getByEmail(email)
		if existingUser == nil {
			return errors.Errorf("User Not Found.")
		}

		uuidBinary, err := uuid.Must(uuid.Parse(existingUser.Uuid)).MarshalBinary()
		if err != nil {
			return err
		}

		token, err := tokenCreate(uuidBinary)
		if err != nil {
			return err
		}

		data := map[string]interface{}{
			"host":  "http://localhost:4000", // @todo this needs to be configurable
			"token": base64.RawURLEncoding.EncodeToString(token),
		}

		buf := &bytes.Buffer{}
		err = tplEmailText.Execute(buf, data)
		if err != nil {
			return err
		}
		textMsg := buf.String()

		buf.Reset()
		err = tplEmailHtml.Execute(buf, data)
		if err != nil {
			return err
		}
		htmlMsg := buf.String()

		fmt.Println(textMsg)
		fmt.Println(htmlMsg)
		//err = sendEmail(existingUser.Name, existingUser.Email, textMsg, htmlMsg)
		//if err != nil {
		//	return err
		//}

		return nil

	}

	return func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse: "+err.Error(), http.StatusInternalServerError)
			return
		}
		email := r.FormValue("email")
		err = sendResetLink(email)

		if err != nil {
			http.Error(w, "Unable to parse: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tplSuccess.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}

}
