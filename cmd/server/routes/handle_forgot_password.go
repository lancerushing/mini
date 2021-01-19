package routes

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/pkg/errors"
)

func (s *Server) handleForgotPasswordForm() http.HandlerFunc {

	tpl := s.mustSetupTemplate("cmd/server/routes/templates/forgotPasswordForm.html")

	return func(w http.ResponseWriter, r *http.Request) {

		data := map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
		}
		_ = tpl.Execute(w, data)

	}

}

// ################### Submit ###################

func (s *Server) handleForgotPasswordSubmit() http.HandlerFunc {
	tplSuccess := s.mustSetupTemplate("cmd/server/routes/templates/forgotPasswordSuccess.html")

	tplEmailHTML := template.Must(template.ParseFiles("cmd/server/routes/templates/forgotPasswordEmail.html"))
	tplEmailText := template.Must(template.ParseFiles("cmd/server/routes/templates/forgotPasswordEmail.text"))

	sendResetLink := func(email string) error {
		if len(email) == 0 {
			return errors.Errorf("Email is empty")
		}

		existingUser := s.getByEmail(email)
		if existingUser == nil {
			return errors.Errorf("User Not Found.")
		}

		uuidBinary, err := uuid.Must(uuid.Parse(existingUser.UUID)).MarshalBinary()
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
		err = tplEmailHTML.Execute(buf, data)
		if err != nil {
			return err
		}
		htmlMsg := buf.String()

		fmt.Println(textMsg)
		fmt.Println(htmlMsg)

		err = sendEmail(existingUser.Name, existingUser.Email, textMsg, htmlMsg)
		if err != nil {
			return err
		}

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
