package apiserver

import (
	"database/sql"
	"github.com/nozimy/technopark-db-forum/internal/store/create"
	"log"
	"net/http"
	"time"
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

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Hour)
	if err := create.CreateTables(db); err != nil {
		return nil, err
	}

	return db, nil
}
