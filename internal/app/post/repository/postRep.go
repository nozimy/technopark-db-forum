package postRepository

import (
	"database/sql"
	"strconv"
	"technopark-db-forum/internal/app/post"
	"technopark-db-forum/internal/model"
)

type PostRepository struct {
	db *sql.DB
}

func (p PostRepository) Update(id string, message string) (*model.Post, error) {
	postObj := &model.Post{}

	id2, _ := strconv.Atoi(id)

	if err := p.db.QueryRow(
		"UPDATE posts SET message = $2, isEdited = TRUE WHERE id = $1 RETURNING author, created, forum, id, message, thread, isEdited",
		id2,
		message,
	).Scan(
		&postObj.Author,
		&postObj.Created,
		&postObj.Forum,
		&postObj.ID,
		&postObj.Message,
		&postObj.Thread,
		&postObj.IsEdited,
	); err != nil {
		return nil, err
	}

	return postObj, nil
}

func (p PostRepository) FindById(id string) (*model.PostFull, error) {
	postObj := &model.PostFull{}
	postObj.Post = &model.Post{}

	id2, _ := strconv.Atoi(id)

	if err := p.db.QueryRow(
		"SELECT author, created, forum, id, message, thread, isedited FROM posts WHERE id = $1",
		id2,
	).Scan(
		&postObj.Post.Author,
		&postObj.Post.Created,
		&postObj.Post.Forum,
		&postObj.Post.ID,
		&postObj.Post.Message,
		&postObj.Post.Thread,
		&postObj.Post.IsEdited,
	); err != nil {
		return nil, err
	}

	return postObj, nil
}

func NewPostRepository(db *sql.DB) post.Repository {
	return &PostRepository{db}
}
