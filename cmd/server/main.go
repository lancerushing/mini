package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/lancerushing/mini/cmd/server/routes"
	"github.com/lancerushing/mini/lib/logutil"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Run Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	logutil.Configure()

	// @todo How to handle configs?  Yaml? ENV??
	viper.SetConfigName("mini")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	var config routes.Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	s, err := routes.NewServer(&config)
	if err != nil {
		return err
	}

	addr := getAddress()

	log.Info().Msg("starting server on: http://" + addr)
	// @todo Make the keys configurable
	CSRF := csrf.Protect([]byte("01234567890123456789012345678912"), csrf.Secure(false), csrf.Path("/"))

	return http.ListenAndServe(addr, CSRF(s.GetHandler()))
}

func getAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	return "127.0.0.1:" + port
}
