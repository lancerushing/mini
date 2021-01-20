package routes

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleSignupForm() http.HandlerFunc {
	tpl := s.mustSetupTemplate("templates/signupForm.html")

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

func (s *Server) handleSignupSubmit() http.HandlerFunc {
	tplSuccess := s.mustSetupTemplate("templates/signupSuccess.html")

	saveUser := func(userDto *UserDto) error {
		sql := `INSERT INTO users (uuid, name, email, password) VALUES (:uuid, :name, :email, :password)`

		_, err := s.db.NamedExec(sql, &userDto)

		return err
	}

	createUser := func(name string, email string, password string) (*UserDto, error) {
		if len(name) == 0 {
			return nil, errors.Errorf("Name is empty")
		}

		if len(email) == 0 {
			return nil, errors.Errorf("Email is empty")
		}

		if len(password) == 0 {
			return nil, errors.Errorf("Password is empty")
		}

		existingUser := s.getByEmail(email)
		if existingUser != nil {
			return nil, errors.Errorf("Email already taken.")
		}

		bcryptBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

		user := &UserDto{
			UUID:     uuid.New().String(),
			Name:     name,
			Email:    email,
			Password: string(bcryptBytes),
		}

		err := saveUser(user)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse: "+err.Error(), http.StatusInternalServerError)

			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		userDto, err := createUser(name, email, password)
		if err != nil {
			http.Error(w, "Unable to create: "+err.Error(), http.StatusInternalServerError)

			return
		}

		data := map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
			"name":           userDto.Name,
			"email":          userDto.Email,
		}

		err = tplSuccess.Execute(w, data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
