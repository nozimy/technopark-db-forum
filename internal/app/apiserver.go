package apiserver

import (
	"database/sql"
	"log"
	"net/http"
	"technopark-db-forum/internal/store/create"
)

func Start() error {
	config := NewConfig()

	server, err := NewServer(config)
	if err != nil {
		return err
	}

	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println(err)
		}
	}()

	server.ConfigureServer(db)

	return http.ListenAndServe(config.BindAddr, server)
}

func newDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	if err := create.CreateTables(db); err != nil {
		return nil, err
	}

	return db, nil
}
