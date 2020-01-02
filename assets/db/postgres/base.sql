DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS forums;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS votes;


CREATE TABLE IF NOT EXISTS users
(
    id       bigserial not null primary key,
    nickname varchar   not null unique,
    about    varchar,
    email    varchar   not null,
    fullname varchar   not null
);

CREATE UNIQUE INDEX IF NOT EXISTS users_email_uindex
    ON users (LOWER(email));

CREATE TABLE IF NOT EXISTS forums
(
    id       bigserial not null primary key,
    slug     varchar   not null,
    userNick varchar   not null,
    title    varchar,
    posts    int default 0,
    threads  int default 0
);

CREATE UNIQUE INDEX IF NOT EXISTS forums_slug_uindex
    ON forums (LOWER(slug));

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

CREATE TABLE IF NOT EXISTS posts
(
    id       bigserial not null primary key,
    parent   bigint             DEFAULT NULL,
    thread   int,
    forum    varchar,
    author   varchar,
    created  timestamptz        DEFAULT now(),
    isEdited bool               DEFAULT FALSE,
    message  text,
    path     bigint[]  NOT NULL DEFAULT '{0}'
);
-- CREATE INDEX idx_path ON posts (path);
CREATE INDEX IF NOT EXISTS idx_path ON posts USING GIN (path);

CREATE TABLE IF NOT EXISTS votes
(
    nickname varchar  NOT NULL,
    thread   int      NOT NULL,
    voice    smallint NOT NULL CHECK (voice = 1 OR voice = -1),
    PRIMARY KEY (nickname, thread)
)