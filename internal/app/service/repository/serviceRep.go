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
	if _, err := s.db.Exec("TRUNCATE votes, users, posts, threads, forums RESTART IDENTITY CASCADE"); err != nil {
		return err
	}

	return nil
}

func (s ServiceRepository) GetStatus() (*model.Status, error) {
	status := &model.Status{}

	if err := s.db.QueryRow("SELECT "+
		"(SELECT count(*) from forums) AS forum, "+
		"(SELECT count(*) from posts) AS post, "+
		"(SELECT count(*) from threads) AS thread, "+
		"(SELECT count(*) from users) AS user",
	).Scan(&status.Forum, &status.Post, &status.Thread, &status.User); err != nil {
		return nil, err
	}

	return status, nil
}

func NewServiceRepository(db *sql.DB) service.Repository {
	return &ServiceRepository{db}
}
