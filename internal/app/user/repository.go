package user

import "github.com/nozimy/technopark-db-forum/internal/model"

type Repository interface {
	Create(forum *model.User) error
	FindByNickname(nickname string) (*model.User, error)
	Find(nickname string, email string) ([]model.User, error)
	Update(user *model.User) (*model.User, error)
}
