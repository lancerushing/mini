package server

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/matryer/is"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestHandleLoginForm(t *testing.T) {
	srv := setup(t)

	check := is.New(t)

	req, err := http.NewRequest("GET", "/user/login/", nil)
	check.NoErr(err)

	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	check.Equal(w.Code, http.StatusOK)

	check.True(strings.Contains(w.Body.String(), "Login"))

}

func TestHandleLoginSubmit_NoInput(t *testing.T) {
	srv := setup(t)

	check := is.New(t)

	req, err := http.NewRequest("POST", "/user/login/", nil)
	check.NoErr(err)

	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	check.Equal(w.Code, http.StatusBadRequest)

	check.True(strings.Contains(w.Body.String(), "missing form body"))

}
func TestHandleLoginSubmit_EmptyInput(t *testing.T) {
	srv := setup(t)

	check := is.New(t)

	data := url.Values{}
	data.Set("email", "")

	req, err := http.NewRequest("POST", "/user/login/", strings.NewReader(data.Encode()))
	check.NoErr(err)

	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	check.Equal(w.Code, http.StatusBadRequest)

	actualBody := w.Body.String()

	check.True(strings.Contains(actualBody, "Email is empty"))
	check.True(strings.Contains(actualBody, "Password is empty"))

}

func TestHandleLoginSubmit_BadInput(t *testing.T) {
	srv, mock := setupWithMock(t)

	columns := []string{"uuid", "email", "password"}

	mock.ExpectQuery("SELECT uuid, email, password FROM users WHERE email = (.+)").
		WillReturnRows(sqlmock.NewRows(columns))

	check := is.New(t)

	data := url.Values{}
	data.Set("email", "unknown@email.com")
	data.Set("password", "foo")

	req, err := http.NewRequest("POST", "/user/login/", strings.NewReader(data.Encode()))
	check.NoErr(err)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	check.Equal(w.Code, http.StatusBadRequest)

	actualBody := w.Body.String()

	ok := strings.Contains(actualBody, "Email not found")
	if !ok {
		t.Log(actualBody)
	}
	check.True(ok)

}

func setupWithMock(t *testing.T) (*Server, sqlmock.Sqlmock) {
	testSrv := Server{}

	testSrv.logger, _ = zap.NewDevelopment()
	testSrv.layout = template.Must(template.New("test_layout").Parse(`{{ block "main" . }}test layout main{{ end }}s`))
	testSrv.routes()
	testSrv.loginAuth = newAuth("test-auth", "abcdefghijklmnopqrstuvwx", "abcdefghijklmnopqrstuvwx")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer db.Close()

	testSrv.db = sqlx.NewDb(db, "postgres")

	return &testSrv, mock
}

func TestHandleLoginSubmit_BadPasswordInput(t *testing.T) {
	srv, mock := setupWithMock(t)

	columns := []string{"uuid", "email", "password"}

	mock.ExpectQuery("SELECT uuid, email, password FROM users WHERE email = (.+)").
		WithArgs("known@email.com").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("fake-uuid,known@email.com,badPassword"))

	check := is.New(t)

	data := url.Values{}
	data.Set("email", "known@email.com")
	data.Set("password", "foo")

	req, err := http.NewRequest("POST", "/user/login/", strings.NewReader(data.Encode()))
	check.NoErr(err)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	check.Equal(w.Code, http.StatusBadRequest)

	actualBody := w.Body.String()

	check.True(strings.Contains(actualBody, "Bad Password"))

}

func TestHandleLoginSubmit_GoodInput(t *testing.T) {
	srv, mock := setupWithMock(t)

	columns := []string{"uuid", "email", "password"}

	testPassword := "test"
	bcryptBytes, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.MinCost)

	mock.ExpectQuery("SELECT uuid, email, password FROM users WHERE email = (.+)").
		WithArgs("known@email.com").
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("fake-uuid,known@email.com," + string(bcryptBytes)))

	check := is.New(t)

	data := url.Values{}
	data.Set("email", "known@email.com")
	data.Set("password", testPassword)

	req, err := http.NewRequest("POST", "/user/login/", strings.NewReader(data.Encode()))
	check.NoErr(err)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	check.Equal(w.Code, http.StatusSeeOther)

	actualLocation := w.Header().Get("Location")

	check.Equal(actualLocation, "/user/")

}
