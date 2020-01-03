package serviceRepository

import (
	"database/sql"
	"technopark-db-forum/internal/app/service"
	"technopark-db-forum/internal/model"
)

type ServiceRepository struct {
	db *sql.DB
}

func (s ServiceRepository) ClearAll() error {
	if _, err := s.db.Exec("TRUNCATE TABLE votes"); err != nil {
		return err
	}

	if _, err := s.db.Exec("TRUNCATE TABLE users"); err != nil {
		return err
	}

	if _, err := s.db.Exec("TRUNCATE TABLE posts"); err != nil {
		return err
	}

	if _, err := s.db.Exec("TRUNCATE TABLE threads"); err != nil {
		return err
	}

	if _, err := s.db.Exec("TRUNCATE TABLE forums"); err != nil {
		return err
	}

	return nil
}

func (s ServiceRepository) GetStatus() (*model.Status, error) {
	status := &model.Status{}

	if err := s.db.QueryRow("SELECT count(*) from forums").Scan(&status.Forum); err != nil {
		return nil, err
	}

	if err := s.db.QueryRow("SELECT count(*) from posts").Scan(&status.Post); err != nil {
		return nil, err
	}

	if err := s.db.QueryRow("SELECT count(*) from threads").Scan(&status.Thread); err != nil {
		return nil, err
	}

	if err := s.db.QueryRow("SELECT count(*) from users").Scan(&status.User); err != nil {
		return nil, err
	}

	return status, nil
}

func NewServiceRepository(db *sql.DB) service.Repository {
	return &ServiceRepository{db}
}
