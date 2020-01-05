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
	conditionSign := ">="
	limits := params["limit"]
	var limit string = "100"
	if len(limits) >= 1 {
		limit = limits[0]
	}
	descs := params["desc"]
	var desc string = ""
	if len(descs) >= 1 && descs[0] == "true" {
		desc = "desc"
		conditionSign = "<="
	}
	sinces := params["since"]
	var since string = ""
	if len(sinces) >= 1 {
		since = sinces[0]
	}

	//var threads model.Threads
	threads := model.Threads{}

	var query string

	query = "SELECT id, forum, author, slug, created, title, message, votes FROM threads WHERE LOWER(forum) = LOWER($1) "

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

func (r ForumRepository) FindForumUsers(forumSlug string, params map[string][]string) (model.Users, error) {
	limits := params["limit"]
	var limit string = "100"
	if len(limits) >= 1 {
		limit = limits[0]
	}
	sinceConditionSign := ">"
	descs := params["desc"]
	var desc string = ""
	if len(descs) >= 1 && descs[0] == "true" {
		desc = "desc"
		sinceConditionSign = "<"
	}
	sinces := params["since"]
	var since string = ""
	if len(sinces) >= 1 {
		since = sinces[0]
	}

	users := model.Users{}

	var query string

	if true {
		query = "SELECT nickname, email, fullname, about FROM users " +
			"WHERE nickname IN (" +
			"SELECT DISTINCT author FROM threads WHERE LOWER(forum) = LOWER($1) " +
			"UNION SELECT DISTINCT author FROM posts WHERE LOWER(forum) = LOWER($1) " +
			")"
	}
	if since != "" {
		query += " AND LOWER(nickname COLLATE \"POSIX\") " + sinceConditionSign + " LOWER('" + since + "' COLLATE \"POSIX\") "
	}
	query += " ORDER BY LOWER(nickname COLLATE \"POSIX\") " + desc + " LIMIT " + limit

	rows, err := r.db.Query(query, forumSlug)

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
		"SELECT slug, title, usernick, posts, threads FROM forums WHERE LOWER(slug) = LOWER($1)",
		slug,
	).Scan(
		&f.Slug,
		&f.Title,
		&f.User,
		&f.Posts,
		&f.Threads,
	); err != nil {
		return nil, err
	}

	return f, nil
}

func NewForumRepository(db *sql.DB) forum.Repository {
	return &ForumRepository{db}
}
