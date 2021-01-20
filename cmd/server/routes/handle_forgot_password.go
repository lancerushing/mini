package routes

import (
	"bytes"

	// Need embed for templates.
	_ "embed"
	"encoding/base64"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
)

func (s *Server) handleForgotPasswordForm() http.HandlerFunc {
	tpl := s.mustSetupTemplate("templates/forgotPasswordForm.html")

	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
		}
		_ = tpl.Execute(w, data)
	}
}

// ################### Submit ###################

func (s *Server) handleForgotPasswordSubmit() http.HandlerFunc {
	tplSuccess := s.mustSetupTemplate("templates/forgotPasswordSuccess.html")

	//go:embed templates/forgotPasswordEmail.html
	var emailHTML string

	//go:embed templates/forgotPasswordEmail.text
	var emailText string

	tplEmailHTML := template.Must(template.New("html").Parse(emailHTML))
	tplEmailText := template.Must(template.New("text").Parse(emailText))

	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse: "+err.Error(), http.StatusInternalServerError)

			return
		}

		email := r.FormValue("email")

		existingUser := s.getByEmail(email)
		if existingUser == nil {
			http.Error(w, "no such user", http.StatusInternalServerError)

			return
		}

		err = sendResetLink(*existingUser, tplEmailHTML, tplEmailText)

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

func sendResetLink(existingUser UserDto, tplEmailHTML *template.Template, tplEmailText *template.Template) error {
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

	err = sendEmail(existingUser.Name, existingUser.Email, textMsg, htmlMsg)
	if err != nil {
		return err
	}

	return nil
}
