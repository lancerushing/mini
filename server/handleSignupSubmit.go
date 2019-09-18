package server

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserDto struct {
	Uuid     string
	Email    string
	Name     string
	Password string
}

func (s *server) handleSignupSubmit() http.HandlerFunc {

	tplSuccess := s.mustSetupTemplate("server/templates/signupSuccess.html")

	getByEmail := func(email string) *UserDto {
		result := UserDto{}

		err := s.db.Get(&result, "SELECT uuid, name, email, password FROM users WHERE email = $1", email)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		fmt.Printf("%#v\n", email)
		fmt.Printf("%#v\n", result)

		return &result
	}

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

		existingUser := getByEmail(email)
		fmt.Printf("%#v\n", existingUser)
		if existingUser != nil {
			return nil, errors.Errorf("Email already taken.")
		}

		bcryptBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

		user := &UserDto{
			Uuid:     uuid.New().String(),
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
