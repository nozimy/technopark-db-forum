package post

import "github.com/nozimy/technopark-db-forum/internal/model"

type Usecase interface {
	FindById(id string, params map[string][]string) (*model.PostFull, error)
	Update(id string, message string) (*model.Post, error)
}
