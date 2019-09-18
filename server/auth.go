package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/securecookie"
)

type auth struct {
	cookieName string
	ctxKey     string
	sc         *securecookie.SecureCookie
}

func NewAuth(name string, cookieHashKey string, cookieEncryptKey string) *auth {
	a := &auth{}
	a.cookieName = name
	a.ctxKey = name
	a.sc = securecookie.New([]byte(cookieHashKey), []byte(cookieEncryptKey))

	return a
}

func (a *auth) setCookie(w http.ResponseWriter, uuid string) error {

	value := map[string]string{
		"uuid": uuid,
	}

	encoded, err := a.sc.Encode(a.cookieName, value)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     a.cookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)

	return nil
}

func (a *auth) deleteCooke(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     a.cookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
}

func (a *auth) middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie(a.cookieName)

		if err != nil { // no cookie
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		value := make(map[string]string)
		if err = a.sc.Decode(a.cookieName, cookie.Value, &value); err != nil { // bad cookie
			a.deleteCooke(w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if len(value["uuid"]) == 0 { // empty value
			a.deleteCooke(w)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// put it on the context
		ctx := r.Context()
		ctx = context.WithValue(ctx, a.ctxKey, value["uuid"])

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
