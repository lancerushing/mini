package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) routes() {
	s.router = chi.NewRouter()

	// middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.addAnonRoutes()

	s.router.Group(s.authRoutes)
	s.router.Group(s.resetPasswordRoutes)
}

func (s *Server) addAnonRoutes() {
	mux := s.router

	mux.Get("/", s.handleIndex())

	mux.Get("/user/signup/", s.handleSignupForm())
	mux.Post("/user/signup/", s.handleSignupSubmit())

	mux.Get("/user/login/", s.handleLoginForm())
	mux.Post("/user/login/", s.handleLoginSubmit())

	mux.Get("/user/forgot-password/", s.handleForgotPasswordForm())
	mux.Post("/user/forgot-password/", s.handleForgotPasswordSubmit())

	mux.Get("/user/reset-password/verify/", s.handleResetPasswordVerify())
}

func (s *Server) resetPasswordRoutes(mux chi.Router) {

	mux.Use(s.pwResetAuth.middleware)

	mux.Get("/user/reset-password/", s.handleRestPasswordForm())
	mux.Post("/user/reset-password/", s.handleRestPasswordSubmit())

}

func (s *Server) authRoutes(mux chi.Router) {

	mux.Use(s.loginAuth.middleware)

	mux.Get("/user/", s.handleUserDetails())
	mux.Post("/user/logout/", s.handleLogout())

}
