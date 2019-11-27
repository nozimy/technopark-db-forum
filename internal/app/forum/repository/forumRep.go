package forumRepository

import (
	"database/sql"
	"technopark-db-forum/internal/app/forum"
	"technopark-db-forum/internal/model"
)

type ForumRepository struct {
	db *sql.DB
}

func (r ForumRepository) Create(f *model.Forum) error {
	return r.db.QueryRow(
		"INSERT INTO forums (name) "+
			"VALUES ($1) RETURNING id",
		f.Name,
	).Scan(&f.ID)
}

func (r ForumRepository) Find(id int64) (*model.Forum, error) {
	u := &model.Forum{}
	if err := r.db.QueryRow(
		"SELECT * FROM forums WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Name,
	); err != nil {
		return nil, err
	}
	return u, nil
}

func NewForumRepository(db *sql.DB) forum.Repository {
	return &ForumRepository{db}
}
