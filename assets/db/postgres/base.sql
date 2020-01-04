DROP INDEX IF EXISTS idx_users_email_uindex;
DROP INDEX IF EXISTS idx_users_nickname_uindex;
DROP INDEX IF EXISTS idx_forums_slug_uindex;
DROP INDEX IF EXISTS idx_forums_userNick_unique;;
DROP INDEX IF EXISTS idx_threads_slug;
DROP INDEX IF EXISTS idx_threads_forum;
DROP INDEX IF EXISTS idx_posts_forum;
DROP INDEX IF EXISTS idx_posts_parent;
DROP INDEX IF EXISTS idx_posts_path;
DROP INDEX IF EXISTS idx_posts_thread;
DROP INDEX IF EXISTS idx_posts_thread_id;

-- DROP TRIGGER IF EXISTS forum_posts_count_ins ON posts;
-- DROP TRIGGER IF EXISTS forum_posts_count_del ON posts;
-- DROP TRIGGER IF EXISTS forum_threads_count_ins ON threads;
-- DROP TRIGGER IF EXISTS forum_threads_count_del ON threads;
DROP TRIGGER IF EXISTS on_vote_insert ON votes;
DROP TRIGGER IF EXISTS on_vote_update ON votes;

-- DROP FUNCTION IF EXISTS fn_update_threads_count();
-- DROP FUNCTION IF EXISTS fn_update_posts_count();
DROP FUNCTION IF EXISTS fn_update_thread_votes_ins();
DROP FUNCTION IF EXISTS fn_update_thread_votes_upd();

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS votes;


CREATE TABLE IF NOT EXISTS users
(
    id       bigserial not null primary key,
    nickname varchar(50)   not null unique,
    about    text,
    email    varchar(50)   not null,
    fullname varchar(100)   not null
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_uindex
    ON users (LOWER(email));
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_nickname_uindex
    ON users (LOWER(nickname));

CREATE TABLE IF NOT EXISTS forums
(
    id       bigserial not null primary key,
    slug     varchar   not null,
    userNick varchar   not null,
    title    varchar,
    posts    int default 0,
    threads  int default 0
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_forums_slug_uindex
    ON forums (LOWER(slug));
CREATE UNIQUE INDEX IF NOT EXISTS idx_forums_userNick_unique
    ON forums (LOWER(userNick));

CREATE TABLE IF NOT EXISTS threads
(
    id      serial not null primary key,
    slug    varchar,
    title   varchar,
    message varchar,
    votes   int         default 0,
    author  varchar,
    forum   varchar,
    created timestamptz DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_threads_slug
    ON threads (LOWER(slug));
CREATE INDEX IF NOT EXISTS idx_threads_forum
    ON threads (LOWER(forum));

CREATE TABLE IF NOT EXISTS posts
(
    id       bigserial not null primary key,
    parent   bigint             DEFAULT NULL,
    path     bigint[]  NOT NULL DEFAULT '{0}',
    thread   int REFERENCES threads(id) NOT NULL,
    forum    varchar,
    author   varchar,
    created  timestamptz        DEFAULT now(),
    isEdited bool               DEFAULT FALSE,
    message  text
);
CREATE INDEX IF NOT EXISTS idx_posts_path ON posts USING GIN (path);
CREATE INDEX IF NOT EXISTS idx_posts_thread ON posts (thread);
CREATE INDEX IF NOT EXISTS idx_posts_forum ON posts (forum);
CREATE INDEX IF NOT EXISTS idx_posts_parent ON posts (parent);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id ON posts (thread, id);

CREATE TABLE IF NOT EXISTS votes
(
    nickname varchar  REFERENCES users(nickname) NOT NULL,
    thread   int      REFERENCES threads(id) NOT NULL,
    voice    smallint NOT NULL CHECK (voice = 1 OR voice = -1),
    PRIMARY KEY (nickname, thread)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_nickname_thread_unique
    ON votes (LOWER(nickname), thread);

CREATE FUNCTION fn_update_thread_votes_ins()
    RETURNS TRIGGER AS '
    BEGIN
        UPDATE threads
        SET
            votes = votes + NEW.voice
        WHERE id = NEW.thread;
        RETURN NULL;
    END;
' LANGUAGE plpgsql;


CREATE TRIGGER on_vote_insert
    AFTER INSERT ON votes
    FOR EACH ROW EXECUTE PROCEDURE fn_update_thread_votes_ins();

CREATE FUNCTION fn_update_thread_votes_upd()
    RETURNS TRIGGER AS '
    BEGIN
        IF OLD.voice = NEW.voice
        THEN
            RETURN NULL;
        END IF;
        UPDATE threads
        SET
            votes = votes + CASE WHEN NEW.voice = -1
                                     THEN -2
                                 ELSE 2 END
        WHERE id = NEW.thread;
        RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_vote_update
    AFTER UPDATE ON votes
    FOR EACH ROW EXECUTE PROCEDURE fn_update_thread_votes_upd();


-- CREATE OR REPLACE FUNCTION fn_update_posts_count() RETURNS trigger AS $update_posts_count$
--     BEGIN
--         UPDATE forums AS f SET posts = (SELECT count(id) FROM posts WHERE forum = f.slug)
--         WHERE slug = (SELECT forum FROM new_table LIMIT 1);
--         RETURN NULL;
--     END;
-- $update_posts_count$ LANGUAGE plpgsql;
--
-- CREATE OR REPLACE FUNCTION fn_update_threads_count() RETURNS trigger AS $body$
--     BEGIN
--         UPDATE forums AS f SET threads = (SELECT count(id) FROM threads WHERE forum = f.slug)
--         WHERE slug = (SELECT forum FROM new_table LIMIT 1);
--         RETURN NULL;
--     END;
-- $body$ LANGUAGE plpgsql;

-- CREATE TRIGGER forum_posts_count_ins
--     AFTER INSERT
--     ON posts
--     REFERENCING NEW TABLE AS new_table
--     FOR EACH STATEMENT
--     EXECUTE PROCEDURE fn_update_posts_count();
--
-- CREATE TRIGGER forum_posts_count_del
--     AFTER DELETE
--     ON posts
--     REFERENCING OLD TABLE AS new_table
--     FOR EACH STATEMENT
--     EXECUTE PROCEDURE fn_update_posts_count();

-- CREATE TRIGGER forum_threads_count_ins
--     AFTER INSERT
--     ON threads
--     REFERENCING NEW TABLE AS new_table
--     FOR EACH STATEMENT
--     EXECUTE PROCEDURE fn_update_threads_count();
--
-- CREATE TRIGGER forum_threads_count_del
--     AFTER DELETE
--     ON threads
--     REFERENCING OLD TABLE AS new_table
--     FOR EACH STATEMENT
--     EXECUTE PROCEDURE fn_update_threads_count();