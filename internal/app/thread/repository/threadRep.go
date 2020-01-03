package threadRepository

import (
	"database/sql"
	"github.com/lib/pq"
	"github.com/nozimy/technopark-db-forum/internal/app/thread"
	"github.com/nozimy/technopark-db-forum/internal/model"
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

	//stmt, err := tx.Prepare("INSERT INTO votes(nickname, voice, thread) VALUES ($1, $2, $3) ON CONFLICT ON CONSTRAINT votes_pkey DO UPDATE SET voice = $2")
	stmt, err := tx.Prepare("UPDATE votes SET voice = $3 WHERE LOWER(nickname) = LOWER($1) AND thread=$2")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	stmtThread, err := tx.Prepare("UPDATE threads SET votes = (SELECT sum(voice) from votes WHERE thread = $1) WHERE id=$1 RETURNING slug, title, message, forum, author, created, votes, id")
	if err != nil {
		return nil, err
	}
	defer stmtThread.Close()

	row := tx.QueryRow("SELECT voice FROM votes WHERE LOWER(nickname) = LOWER($1) AND thread=$2", vote.Nickname, thread.ID)
	var prevVoice int32
	if err := row.Scan(&prevVoice); err != nil {
		_, err = tx.Exec("INSERT INTO votes(nickname, thread, voice) VALUES ($1, $2, $3::smallint)", vote.Nickname, thread.ID, vote.Voice)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		_, err = stmt.Exec(vote.Nickname, thread.ID, vote.Voice)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	rowT := stmtThread.QueryRow(thread.ID)
	err = rowT.Scan(
		&thread.Slug,
		&thread.Title,
		&thread.Message,
		&thread.Forum,
		&thread.Author,
		&thread.Created,
		&thread.Votes,
		&thread.ID,
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
	//desc = "asc"

	if sort == "flat" {
		query = "SELECT id, parent, thread, forum, author, created, message, isedited, path FROM posts WHERE thread = $1 "
		if since != "" {
			query += " AND id " + conditionSign + " '" + since + "' "
		}
		query += " ORDER BY created " + desc + ", id " + desc + " LIMIT " + limit
	} else if sort == "tree" {
		//query = "SELECT id, parent, thread, forum, author, created, message, isedited, path FROM posts WHERE thread = $1 AND (array_length(path, 1) IS NULL OR array_length(path, 1) > 0) "
		//if since != "" {
		//	query += " AND id " + conditionSign+ " '" + since + "' "
		//}
		//query += " ORDER BY path[1] " + desc +  ", path " +desc+" LIMIT " + limit
		orderString := " ORDER BY path[1] " + desc + ", path " + desc
		query = "WITH temp as (SELECT id, parent, thread, forum, author, created, message, isedited, path, row_number() " +
			"over (" + orderString + ") as rownum " +
			"FROM posts WHERE thread = $1 AND (array_length(path, 1) IS NULL OR array_length(path, 1) > 0) " +
			orderString + ") " +
			"SELECT id, parent, thread, forum, author, created, message, isedited, path " +
			"FROM temp " +
			"LIMIT " + limit
		if since != "" {
			query += " offset (select rownum from temp where id = " + since + "); "
		}
	} else if sort == "parent_tree" {
		//query = "SELECT id, parent, thread, forum, author, created, message, isedited, path FROM posts WHERE thread = $1 AND array_length(path, 1) = 1 "
		//if since != "" {
		//	query += " AND id " + conditionSign+ " '" + since + "' "
		//}
		//query += " ORDER BY path[1] " + desc + ", path "+desc+" LIMIT " + limit
		query = "WITH temp as (SELECT id, parent, thread, forum, author, created, message, isedited, path, " +
			"row_number() over (ORDER BY path[1] " + desc + ", path) as rownum " +
			"FROM posts " +
			"WHERE thread = $1 AND " +
			"path && (SELECT ARRAY (select id from posts WHERE thread = $1 AND array_length(path, 1) = 1 " +
			"ORDER BY path[1] " + desc + ", path "

		if since == "" {
			query += " LIMIT " + limit
		}

		query += ")) " +
			"ORDER BY path[1] " + desc + ", path) " +
			"SELECT id, parent, thread, forum, author, created, message, isedited, path " +
			"FROM temp " +
			"ORDER BY path[1] " + desc + ", path "

		if since != "" {
			query += " offset (select rownum from temp where id = " + since + "); "
		}
	}

	rows, err := t.db.Query(query, thread.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := model.Post{}
		err := rows.Scan(&p.ID, &p.Parent, &p.Thread, &p.Forum, &p.Author, &p.Created, &p.Message, &p.IsEdited, pq.Array(&p.Path))
		if err != nil {
			return nil, err
		}

		//if sort == "parent_tree" {
		//	childRows, err := t.db.Query("SELECT id, parent, thread, forum, author, created, message, isedited, path FROM posts " +
		//		"WHERE thread = $1 AND path && ARRAY[$2::bigint] ORDER BY path", thread.ID, p.ID)
		//	if err != nil {
		//		return nil, err
		//	}
		//
		//	for childRows.Next() {
		//		p := model.Post{}
		//		err := childRows.Scan(&p.ID, &p.Parent, &p.Thread, &p.Forum, &p.Author, &p.Created, &p.Message, &p.IsEdited, pq.Array(&p.Path))
		//		if err != nil {
		//			return nil, err
		//		}
		//
		//		posts = append(posts, &p)
		//	}
		//	if err := childRows.Close(); err != nil {
		//		return nil, err
		//	}
		//} else {
		posts = append(posts, &p)
		//}
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

	stmt, err := tx.Prepare("INSERT INTO posts(id, parent, thread, forum, author, created, message, path) SELECT nextval('posts_id_seq'::regclass), $1, $2, $3, $4, $5, $6, ARRAY[currval(pg_get_serial_sequence('posts', 'id'))::bigint] RETURNING id, parent, thread, forum, author, created, message, isedited")
	if err != nil {
		return nil, err
	}

	//stmtWithParentCheck, err := tx.Prepare("INSERT INTO posts(parent, thread, forum, author, created, message) SELECT $1, $2, $3, $4, $5, $6 WHERE EXISTS(SELECT id FROM posts WHERE id = $1) RETURNING id, parent, thread, forum, author, created, message, isedited")
	stmtWithParentCheck, err := tx.Prepare("INSERT INTO posts(id, parent, thread, forum, author, created, message, path) SELECT nextval('posts_id_seq'::regclass), $1, $2, $3, $4, $5, $6, path || (select currval(pg_get_serial_sequence('posts', 'id'))::bigint) FROM posts WHERE id = $1 AND thread = $2 RETURNING id, parent, thread, forum, author, created, message, isedited")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	defer stmtWithParentCheck.Close()

	now := time.Now()

	for _, post := range *posts {
		var row *sql.Row
		if post.Parent == 0 {
			row = stmt.QueryRow(post.Parent, thread.ID, thread.Forum, post.Author, now, post.Message)
		} else {
			row = stmtWithParentCheck.QueryRow(post.Parent, thread.ID, thread.Forum, post.Author, now, post.Message)
		}

		err := row.Scan(
			&post.ID,
			&post.Parent,
			&post.Thread,
			&post.Forum,
			&post.Author,
			&post.Created,
			&post.Message,
			&post.IsEdited,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
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

	return th, nil
}

func (t ThreadRepository) FindByIdOrSlug(id int, slug string) (*model.Thread, error) {
	th := &model.Thread{}

	err := t.db.QueryRow(
		"SELECT slug, title, message, forum, author, created, votes, id FROM threads WHERE id=$1 OR LOWER(slug)=LOWER($2) AND slug <> ''",
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
