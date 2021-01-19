package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/spf13/viper"

	"github.com/lancerushing/mini/cmd/server/routes"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Run Error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {

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

	// @todo Make the keys configurable
	CSRF := csrf.Protect([]byte("01234567890123456789012345678912"), csrf.Secure(false), csrf.Path("/"))
	return http.ListenAndServe(":4001", CSRF(s.GetHandler()))

}
