package server

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) handleLoginSubmit() http.HandlerFunc {

	getByEmail := func(email string) (*UserDto, error) {
		result := []UserDto{}

		err := s.db.Select(&result, "SELECT uuid, email, password FROM users WHERE Email = $1 ", email)
		if err != nil {
			return nil, err
		}

		if len(result) != 1 {
			return nil, errors.New("Email not found")
		}
		return &result[0], nil

	}

	return func(w http.ResponseWriter, r *http.Request) {
		var user *UserDto
		var err error

		err = r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse: "+err.Error(), http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		fieldErrors := map[string]interface{}{}

		if len(email) == 0 {
			fieldErrors["emailError"] = "Email is empty"
		}
		if len(password) == 0 {
			fieldErrors["passwordError"] = "Password is empty"
		}

		if len(fieldErrors) == 0 {

			userMatch, err := getByEmail(email)
			if err != nil {
				fieldErrors["emailError"] = err.Error()
			}

			if userMatch != nil {
				err = bcrypt.CompareHashAndPassword([]byte(userMatch.Password), []byte(password))
				if err != nil {
					fieldErrors["passwordError"] = "Bad Password"
				} else {
					user = userMatch
				}

			}

		}

		if user == nil {

			tpl := s.mustSetupTemplate("server/templates/loginForm.html")
			data := map[string]interface{}{
				csrf.TemplateTag: csrf.TemplateField(r),
				"errorMsg":       template.HTML(`<div class="text-danger">Invalid Login</div>`),
				"email":          email,
				"password":       password,
			}
			for k, v := range fieldErrors {
				data[k] = v
			}

			w.WriteHeader(http.StatusBadRequest)
			err = tpl.Execute(w, data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		s.loginAuth.setCookie(w, user.Uuid)
		http.Redirect(w, r, "/user/", http.StatusSeeOther)

	}

}
