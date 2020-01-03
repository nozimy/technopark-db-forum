package user

import "github.com/nozimy/technopark-db-forum/internal/model"

type Usecase interface {
	CreateUser(user *model.User) ([]model.User, error)
	Find(nickname string) (*model.User, error)
	Update(user *model.User) (*model.User, error, int)
}
