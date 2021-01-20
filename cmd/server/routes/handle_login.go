package routes

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleLoginForm() http.HandlerFunc {
	tpl := s.mustSetupTemplate("templates/loginForm.html")

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

func (s *Server) handleLoginSubmit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user *UserDto
			err  error
		)

		err = r.ParseForm()
		if err != nil {
			http.Error(w, "unable to parse: "+err.Error(), http.StatusBadRequest)

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
			userMatch := s.getByEmail(email)
			if userMatch == nil {
				fieldErrors["emailError"] = "Email not found"
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
			tpl := s.mustSetupTemplate("templates/loginForm.html")
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

		err = s.loginAuth.setCookie(w, user.UUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		http.Redirect(w, r, "/user/", http.StatusSeeOther)
	}
}

// ################### Logout ###################

func (s *Server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.loginAuth.deleteCooke(w)
		http.Redirect(w, r, "../", http.StatusSeeOther)
	}
}
