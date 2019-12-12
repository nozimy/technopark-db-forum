package postRepository

import (
	"database/sql"
	"technopark-db-forum/internal/app/post"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) post.Repository {
	return &PostRepository{db}
}
