package server

import (
	"net/http"
)

func (s *Server) handleLogout() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		s.loginAuth.deleteCooke(w)
		http.Redirect(w, r, "../", http.StatusSeeOther)
	}

}
