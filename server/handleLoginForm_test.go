package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestHandleLoginForm(t *testing.T) {
	srv := setup(t)

	is := is.New(t)

	req, err := http.NewRequest("GET", "/user/login/", nil)
	is.NoErr(err)

	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	is.Equal(w.Code, http.StatusOK)

	is.True(strings.Contains(w.Body.String(), "Login"))

}
