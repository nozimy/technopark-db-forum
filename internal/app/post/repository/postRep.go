package postRepository

import (
	"database/sql"
	"github.com/nozimy/technopark-db-forum/internal/app/post"
	"github.com/nozimy/technopark-db-forum/internal/model"
	"strconv"
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

func (p PostRepository) FindById(id string, includeUser, includeForum, includeThread bool) (*model.PostFull, error) {
	postObj := &model.PostFull{}
	postObj.Post = &model.Post{}

	id2, _ := strconv.Atoi(id)

	if err := p.db.QueryRow(
		"SELECT author, created, forum, id, message, thread, isedited, parent FROM posts WHERE id = $1",
		id2,
	).Scan(
		&postObj.Post.Author,
		&postObj.Post.Created,
		&postObj.Post.Forum,
		&postObj.Post.ID,
		&postObj.Post.Message,
		&postObj.Post.Thread,
		&postObj.Post.IsEdited,
		&postObj.Post.Parent,
	); err != nil {
		return nil, err
	}

	if includeUser {
		postObj.Author = &model.User{}
		if err := p.db.QueryRow(
			"SELECT about, email, fullname, nickname FROM users WHERE LOWER(nickname) = LOWER($1)",
			postObj.Post.Author,
		).Scan(
			&postObj.Author.About,
			&postObj.Author.Email,
			&postObj.Author.Fullname,
			&postObj.Author.Nickname,
		); err != nil {
			return nil, err
		}
	}

	if includeForum {
		postObj.Forum = &model.Forum{}
		if err := p.db.QueryRow(
			"SELECT usernick, title, slug, posts, threads FROM forums WHERE LOWER(slug) = LOWER($1)",
			postObj.Post.Forum,
		).Scan(
			&postObj.Forum.User,
			&postObj.Forum.Title,
			&postObj.Forum.Slug,
			&postObj.Forum.Posts,
			&postObj.Forum.Threads,
		); err != nil {
			return nil, err
		}
	}

	if includeThread {
		postObj.Thread = &model.Thread{}
		if err := p.db.QueryRow(
			"SELECT forum, slug, title, author, message, id, created, votes FROM threads WHERE id = $1",
			postObj.Post.Thread,
		).Scan(
			&postObj.Thread.Forum,
			&postObj.Thread.Slug,
			&postObj.Thread.Title,
			&postObj.Thread.Author,
			&postObj.Thread.Message,
			&postObj.Thread.ID,
			&postObj.Thread.Created,
			&postObj.Thread.Votes,
		); err != nil {
			return nil, err
		}
	}

	return postObj, nil
}

func NewPostRepository(db *sql.DB) post.Repository {
	return &PostRepository{db}
}
