package apiserver

import (
	"net/http"
)

func Start() error {
	config := NewConfig()

	server, err := NewServer(config)
	if err != nil {
		return err
	}

	server.ConfigureServer()

	return http.ListenAndServe(config.BindAddr, server)
}
