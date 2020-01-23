package userRepository

import (
	"database/sql"
	"github.com/nozimy/technopark-db-forum/internal/app/user"
	"github.com/nozimy/technopark-db-forum/internal/model"
	cache2 "github.com/patrickmn/go-cache"
)

type UserRepository struct {
	db    *sql.DB
	cache *cache2.Cache
}

func (u UserRepository) Create(user *model.User) error {
	return u.db.QueryRow(
		"INSERT INTO users (nickname, email, about, fullname) "+
			"VALUES ($1, $2, $3, $4) RETURNING nickname",
		user.Nickname,
		user.Email,
		user.About,
		user.Fullname,
	).Scan(&user.Nickname)
}

func (u UserRepository) FindByNickname(nickname string) (*model.User, error) {
	userObj := &model.User{}

	if x, found := u.cache.Get(nickname); found {
		userObj = x.(*model.User)
	} else {
		if err := u.db.QueryRow(
			"SELECT nickname, about, email, fullname FROM users WHERE nickname = $1",
			nickname,
		).Scan(
			&userObj.Nickname,
			&userObj.About,
			&userObj.Email,
			&userObj.Fullname,
		); err != nil {
			return nil, err
		}

		u.cache.Set(userObj.Nickname, userObj, cache2.DefaultExpiration)
		u.cache.Set(userObj.Email, userObj, cache2.DefaultExpiration)
	}

	return userObj, nil
}

func (u UserRepository) Find(nickname string, email string) ([]model.User, error) {
	var users []model.User

	if x, found := u.cache.Get(nickname); found {
		userObj := x.(*model.User)
		users = append(users, *userObj)
	} else if x, found := u.cache.Get(nickname); found {
		userObj := x.(*model.User)
		users = append(users, *userObj)
	} else {
		rows, err := u.db.Query(
			"SELECT nickname, about, email, fullname FROM users WHERE nickname = $1 OR email = $2",
			nickname,
			email,
		)

		if err != nil {
			return nil, err
		}

		for rows.Next() {
			obj := model.User{}
			err := rows.Scan(&obj.Nickname, &obj.About, &obj.Email, &obj.Fullname)
			if err != nil {
				return nil, err
			}
			users = append(users, obj)

			u.cache.Set(obj.Nickname, obj, cache2.DefaultExpiration)
			u.cache.Set(obj.Email, obj, cache2.DefaultExpiration)
		}

		if err := rows.Close(); err != nil {
			return nil, err
		}
	}

	return users, nil
}

func (u UserRepository) Update(user *model.User) (*model.User, error) {
	if err := u.db.QueryRow(
		"UPDATE users SET about = $1, email = $2, fullname = $3 WHERE nickname = $4 RETURNING nickname, about, email, fullname",
		user.About,
		user.Email,
		user.Fullname,
		user.Nickname,
	).Scan(
		&user.Nickname,
		&user.About,
		&user.Email,
		&user.Fullname,
	); err != nil {
		return nil, err
	}

	u.cache.Set(user.Nickname, user, cache2.DefaultExpiration)
	u.cache.Set(user.Email, user, cache2.DefaultExpiration)

	return user, nil
}

func NewUserRepository(db *sql.DB, c *cache2.Cache) user.Repository {
	return &UserRepository{db, c}
}
