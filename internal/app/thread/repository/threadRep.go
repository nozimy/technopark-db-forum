package threadRepository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/nozimy/technopark-db-forum/internal/app/thread"
	"github.com/nozimy/technopark-db-forum/internal/model"
	"strconv"
	"strings"
	"time"
)

type ThreadRepository struct {
	db *sql.DB
}

func (t ThreadRepository) Vote(thread *model.Thread, vote *model.Vote) (*model.Thread, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare("INSERT INTO votes(nickname, voice, thread) VALUES ($1, $2, $3) ON CONFLICT ON CONSTRAINT votes_pkey DO UPDATE SET voice = $2")
	//stmt, err := tx.Prepare("UPDATE votes SET voice = $3 WHERE LOWER(nickname) = LOWER($1) AND thread=$2")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	//stmtThread, err := tx.Prepare("UPDATE threads SET votes = (SELECT sum(voice) from votes WHERE thread = $1) WHERE id=$1 RETURNING slug, title, message, forum, author, created, votes, id")
	//if err != nil {
	//	return nil, err
	//}
	//defer stmtThread.Close()

	//row := tx.QueryRow("SELECT voice FROM votes WHERE LOWER(nickname) = LOWER($1) AND thread=$2", vote.Nickname, thread.ID)
	//var prevVoice int32
	//if err := row.Scan(&prevVoice); err != nil {
	//	_, err = tx.Exec("INSERT INTO votes(nickname, thread, voice) VALUES ($1, $2, $3::smallint)", vote.Nickname, thread.ID, vote.Voice)
	//	if err != nil {
	//		tx.Rollback()
	//		return nil, err
	//	}
	//} else {
	_, err = stmt.Exec(vote.Nickname, vote.Voice, thread.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	//}

	//rowT := stmtThread.QueryRow(thread.ID)
	//err = rowT.Scan(
	//	&thread.Slug,
	//	&thread.Title,
	//	&thread.Message,
	//	&thread.Forum,
	//	&thread.Author,
	//	&thread.Created,
	//	&thread.Votes,
	//	&thread.ID,
	//)
	rowT := tx.QueryRow("SELECT votes FROM threads WHERE id = $1", thread.ID)
	err = rowT.Scan(
		&thread.Votes,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (t ThreadRepository) GetThreadPosts(thread *model.Thread, limit, desc, since, sort string) (model.Posts, error) {
	posts := make(model.Posts, 0)

	var query string

	conditionSign := ">"
	if desc == "desc" {
		conditionSign = "<"
	}

	if sort == "flat" {
		query = "SELECT id, parent, thread, forum, author, created, message, isedited FROM posts WHERE thread = $1 "
		if since != "" {
			query += fmt.Sprintf(" AND id %s %s ", conditionSign, since)
		}
		query += fmt.Sprintf(" ORDER BY created %s, id %s LIMIT %s", desc, desc, limit)
	} else if sort == "tree" {
		orderString := fmt.Sprintf(" ORDER BY path[1] %s, path %s ", desc, desc)
		//query = "WITH temp as (" +
		//	"SELECT id, parent, thread, forum, author, created, message, isedited, path, row_number() " +
		//	"over (" + orderString + ") as rownum " +
		//	"FROM posts WHERE thread = $1 AND (array_length(path, 1) IS NULL OR array_length(path, 1) > 0) "
		//if since != "" {
		//	query += fmt.Sprintf(" AND path %s (SELECT path FROM posts WHERE id = %s) ", conditionSign, since)
		//}
		//query+= orderString + ") " +
		//	"SELECT id, parent, thread, forum, author, created, message, isedited, path " +
		//	"FROM temp " +
		//	"LIMIT " + limit

		query = "SELECT id, parent, thread, forum, author, created, message, isedited " +
			"FROM posts " +
			"WHERE thread = $1 "
		if since != "" {
			query += fmt.Sprintf(" AND path %s (SELECT path FROM posts WHERE id = %s) ", conditionSign, since)
		}
		query += orderString
		query += fmt.Sprintf("LIMIT %s", limit)

	} else if sort == "parent_tree" {
		//query = "WITH temp as ( " +
		//	"SELECT id, parent, thread, forum, author, created, message, isedited, path, " +
		//	"row_number() over (ORDER BY path[1] " + desc + ", path) as rownum " +
		//	"FROM posts " +
		//	"WHERE thread = $1 AND " +
		//	//"path && (SELECT ARRAY (select id from posts WHERE thread = $1 AND array_length(path, 1) = 1 " +
		//	"path && (SELECT ARRAY (select id from posts WHERE thread = $1 AND parent = 0 "
		//if since != "" {
		//	query += fmt.Sprintf(" AND path %s (SELECT path[1:1] FROM posts WHERE id = %s) ", conditionSign, since)
		//}
		//query += "ORDER BY path[1] " + desc + ", path "
		//
		//query += " LIMIT " + limit
		//
		//query += ")) " +
		//	"ORDER BY path[1] " + desc + ", path) " +
		//	"SELECT id, parent, thread, forum, author, created, message, isedited, path " +
		//	"FROM temp " +
		//	"ORDER BY path[1] " + desc + ", path "

		query = "SELECT id, parent, thread, forum, author, created, message, isedited " +
			"FROM posts " +
			"WHERE thread = $1 AND path && (SELECT ARRAY (select id from posts WHERE thread = $1 AND parent = 0 "
		if since != "" {
			query += fmt.Sprintf(" AND path %s (SELECT path[1:1] FROM posts WHERE id = %s) ", conditionSign, since)
		}
		query += fmt.Sprintf("ORDER BY path[1] %s, path LIMIT %s)) ", desc, limit)
		query += fmt.Sprintf("ORDER BY path[1] %s, path ", desc)
	}

	rows, err := t.db.Query(query, thread.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := model.Post{}
		err := rows.Scan(&p.ID, &p.Parent, &p.Thread, &p.Forum, &p.Author, &p.Created, &p.Message, &p.IsEdited)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &p)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (t ThreadRepository) UpdateThread(id int, slug string, threadUpdate *model.ThreadUpdate) (*model.Thread, error) {
	th := &model.Thread{}

	err := t.db.QueryRow(
		"UPDATE threads SET title = $1, message = $2 WHERE id=$3 OR LOWER(slug)=LOWER($4) RETURNING slug, title, message, forum, author, created, votes, id",
		threadUpdate.Title,
		threadUpdate.Message,
		id,
		slug,
	).Scan(
		&th.Slug,
		&th.Title,
		&th.Message,
		&th.Forum,
		&th.Author,
		&th.Created,
		&th.Votes,
		&th.ID,
	)

	if err != nil {
		return nil, err
	}

	return th, nil
}

func (t ThreadRepository) CreatePosts(thread *model.Thread, posts *model.Posts) (*model.Posts, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	sqlStr := "INSERT INTO posts(id, parent, thread, forum, author, created, message, path) VALUES "
	vals := []interface{}{}
	for _, post := range *posts {
		var author string
		err = t.db.QueryRow("SELECT nickname FROM users WHERE LOWER(nickname) = LOWER($1)",
			post.Author,
		).Scan(&author)
		if err != nil || author == "" {
			_ = tx.Rollback()
			return nil, errors.New("404")
		}

		if post.Parent == 0 {
			sqlStr += "(nextval('posts_id_seq'::regclass), ?, ?, ?, ?, ?, ?, " +
				"ARRAY[currval(pg_get_serial_sequence('posts', 'id'))::bigint]),"
			vals = append(vals, post.Parent, thread.ID, thread.Forum, post.Author, now, post.Message)
		} else {
			var parentThreadId int32
			err = t.db.QueryRow("SELECT thread FROM posts WHERE id = $1",
				post.Parent,
			).Scan(&parentThreadId)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}
			if parentThreadId != thread.ID {
				_ = tx.Rollback()
				return nil, errors.New("Parent post was created in another thread")
			}

			sqlStr += " (nextval('posts_id_seq'::regclass), ?, ?, ?, ?, ?, ?, " +
				"(SELECT path FROM posts WHERE id = ? AND thread = ?) || " +
				"currval(pg_get_serial_sequence('posts', 'id'))::bigint),"

			vals = append(vals, post.Parent, thread.ID, thread.Forum, post.Author, now, post.Message, post.Parent, thread.ID)
		}

	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	sqlStr += " RETURNING  id, parent, thread, forum, author, created, message, isedited "

	sqlStr = ReplaceSQL(sqlStr, "?")
	if len(*posts) > 0 {
		stmtButch, err := tx.Prepare(sqlStr)
		if err != nil {
			return nil, err
		}
		rows, err := stmtButch.Query(vals...)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		i := 0
		for rows.Next() {
			err := rows.Scan(
				&(*posts)[i].ID,
				&(*posts)[i].Parent,
				&(*posts)[i].Thread,
				&(*posts)[i].Forum,
				&(*posts)[i].Author,
				&(*posts)[i].Created,
				&(*posts)[i].Message,
				&(*posts)[i].IsEdited,
			)
			i += 1

			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	_, err = tx.Exec("UPDATE forums SET posts = posts + $1 WHERE lower(slug) = lower($2)", len(*posts), thread.Forum)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (t ThreadRepository) CreateThread(newThread *model.NewThread) (*model.Thread, error) {
	th := &model.Thread{}
	var row *sql.Row

	if newThread.Created == "" {
		query := "INSERT INTO threads (title, message, forum, author, slug) " +
			"VALUES ($1, $2, $3, $4, $5) RETURNING slug, title, message, forum, author, created, votes, id"
		row = t.db.QueryRow(
			query,
			newThread.Title,
			newThread.Message,
			newThread.Forum,
			newThread.Author,
			newThread.Slug,
		)
	} else {
		query := "INSERT INTO threads (title, message, forum, author, slug, created) " +
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING slug, title, message, forum, author, created, votes, id"
		row = t.db.QueryRow(
			query,
			newThread.Title,
			newThread.Message,
			newThread.Forum,
			newThread.Author,
			newThread.Slug,
			newThread.Created,
		)
	}
	err := row.Scan(
		&th.Slug,
		&th.Title,
		&th.Message,
		&th.Forum,
		&th.Author,
		&th.Created,
		&th.Votes,
		&th.ID,
	)
	if err != nil {
		return nil, err
	}

	_, err = t.db.Exec("UPDATE forums SET threads = threads + 1 WHERE lower(slug) = lower($1)", th.Forum)
	if err != nil {
		return nil, err
	}

	return th, nil
}

func (t ThreadRepository) FindByIdOrSlug(id int, slug string) (*model.Thread, error) {
	th := &model.Thread{}

	err := t.db.QueryRow(
		"SELECT slug, title, message, forum, author, created, votes, id FROM threads WHERE id=$1 OR (LOWER(slug)=LOWER($2) AND slug <> '')",
		id,
		slug,
	).Scan(
		&th.Slug,
		&th.Title,
		&th.Message,
		&th.Forum,
		&th.Author,
		&th.Created,
		&th.Votes,
		&th.ID,
	)

	if err != nil {
		return nil, err
	}

	return th, nil
}

func NewThreadRepository(db *sql.DB) thread.Repository {
	return &ThreadRepository{db}
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}
