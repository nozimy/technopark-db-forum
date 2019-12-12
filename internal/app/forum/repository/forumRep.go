package forumRepository

import (
	"database/sql"
	"log"
	"technopark-db-forum/internal/app/forum"
	"technopark-db-forum/internal/model"
)

type ForumRepository struct {
	db *sql.DB
}

func (r ForumRepository) FindForumThreads(forumSlug string, params map[string][]string) (model.Threads, error) {
	limits := params["limit"]
	var limit string = "100"
	if len(limits) >= 1 {
		limit = limits[0]
	}
	descs := params["desc"]
	var desc string = ""
	if len(descs) >= 1 && descs[0] == "true" {
		desc = "desc"
	}
	sinces := params["since"]
	var since string = ""
	if len(sinces) >= 1 {
		since = sinces[0]
	}

	var threads model.Threads

	var query string

	if true {
		query = "SELECT id, forum, author, slug, created, title, message, votes FROM threads WHERE forum = $1"
	}
	if since != "" {
		query += " AND created > '" + since + "' "
	}
	//query += " ORDER BY created $2 LIMIT $3"
	query += " ORDER BY created " + desc + " LIMIT " + limit

	rows, err := r.db.Query(query, forumSlug)

	if err != nil {
		log.Println("log", err)
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
	descs := params["desc"]
	var desc string = ""
	if len(descs) >= 1 && descs[0] == "true" {
		desc = "desc"
	}
	sinces := params["since"]
	var since string = ""
	if len(sinces) >= 1 {
		since = sinces[0]
	}

	var users model.Users

	var query string

	if true {
		query = "SELECT nickname, email, fullname, about FROM users WHERE nickname IN (SELECT DISTINCT author FROM threads UNION SELECT DISTINCT author FROM posts)"
	}
	if since != "" {
		query += " AND LOWER(nickname) > '" + since + "' "
	}
	query += " ORDER BY LOWER(nickname) " + desc + " LIMIT " + limit

	rows, err := r.db.Query(query)

	if err != nil {
		log.Println(err)
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
		"SELECT slug, title, usernick, posts, threads FROM forums WHERE slug = $1",
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
