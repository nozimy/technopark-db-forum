package forumRepository

import (
	"database/sql"
	"fmt"
	"github.com/nozimy/technopark-db-forum/internal/app/forum"
	"github.com/nozimy/technopark-db-forum/internal/model"
)

type ForumRepository struct {
	db *sql.DB
}

func (r ForumRepository) FindForumThreads(forumSlug string, params map[string][]string) (model.Threads, error) {
	limit := "100"
	if len(params["limit"]) >= 1 {
		limit = params["limit"][0]
	}
	desc := ""
	conditionSign := ">="
	if len(params["desc"]) >= 1 && params["desc"][0] == "true" {
		desc = "desc"
		conditionSign = "<="
	}
	since := ""
	if len(params["since"]) >= 1 {
		since = params["since"][0]
	}

	threads := model.Threads{}

	query := "SELECT id, forum, author, slug, created, title, message, votes FROM threads WHERE LOWER(forum) = LOWER($1) "

	if since != "" {
		query += fmt.Sprintf(" AND created %s '%s' ", conditionSign, since)
	}

	query += fmt.Sprintf(" ORDER BY created %s LIMIT %s", desc, limit)

	rows, err := r.db.Query(query, forumSlug)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := model.Thread{}
		err := rows.Scan(&t.ID, &t.Forum, &t.Author, &t.Slug, &t.Created, &t.Title, &t.Message, &t.Votes)

		if err != nil {
			return nil, err
		}

		threads = append(threads, t)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return threads, nil
}

func (r ForumRepository) FindForumUsers(forumObj *model.Forum, params map[string][]string) (model.Users, error) {
	limit := "100"
	if len(params["limit"]) >= 1 {
		limit = params["limit"][0]
	}
	sinceConditionSign := ">"
	desc := ""
	if len(params["desc"]) >= 1 && params["desc"][0] == "true" {
		desc = "desc"
		sinceConditionSign = "<"
	}
	since := ""
	if len(params["since"]) >= 1 {
		since = params["since"][0]
	}

	users := model.Users{}

	query := "SELECT nickname, email, fullname, about FROM users " +
		"WHERE id IN (SELECT user_id FROM forum_users WHERE forum_id = $1) "
	if since != "" {
		query += fmt.Sprintf(" AND LOWER(nickname) %s LOWER('%s') ", sinceConditionSign, since)
	}
	query += fmt.Sprintf(" ORDER BY LOWER(nickname) %s LIMIT %s ", desc, limit)

	rows, err := r.db.Query(query, forumObj.ID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := model.User{}
		err := rows.Scan(&u.Nickname, &u.Email, &u.Fullname, &u.About)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r ForumRepository) Create(f *model.Forum) error {
	return r.db.QueryRow(
		"INSERT INTO forums (slug, title, usernick) "+
			"VALUES ($1, $2, $3) RETURNING slug, title, usernick, posts, threads",
		f.Slug,
		f.Title,
		f.User,
	).Scan(
		&f.Slug,
		&f.Title,
		&f.User,
		&f.Posts,
		&f.Threads,
	)
}

func (r ForumRepository) Find(slug string) (*model.Forum, error) {
	f := &model.Forum{}

	if err := r.db.QueryRow(
		"SELECT slug, title, usernick, posts, threads, id FROM forums WHERE LOWER(slug) = LOWER($1)",
		slug,
	).Scan(
		&f.Slug,
		&f.Title,
		&f.User,
		&f.Posts,
		&f.Threads,
		&f.ID,
	); err != nil {
		return nil, err
	}

	return f, nil
}

func NewForumRepository(db *sql.DB) forum.Repository {
	return &ForumRepository{db}
}
