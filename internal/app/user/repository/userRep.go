package userRepository

import (
	"database/sql"
	"technopark-db-forum/internal/app/user"
	"technopark-db-forum/internal/model"
)

type UserRepository struct {
	db *sql.DB
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

	if err := u.db.QueryRow(
		"SELECT nickname, about, email, fullname FROM users WHERE LOWER(nickname) = LOWER($1)",
		nickname,
	).Scan(
		&userObj.Nickname,
		&userObj.About,
		&userObj.Email,
		&userObj.Fullname,
	); err != nil {
		return nil, err
	}

	return userObj, nil
}

func (u UserRepository) Find(nickname string, email string) ([]model.User, error) {
	var users []model.User

	rows, err := u.db.Query(
		"SELECT nickname, about, email, fullname FROM users WHERE LOWER(nickname) = LOWER($1) OR LOWER(email) = LOWER($2)",
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
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}

	return users, nil
}

func (u UserRepository) Update(user *model.User) (*model.User, error) {
	if err := u.db.QueryRow(
		"UPDATE users SET about = $1, email = $2, fullname = $3 WHERE LOWER(nickname) = LOWER($4) RETURNING nickname, about, email, fullname",
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

	return user, nil
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &UserRepository{db}
}
