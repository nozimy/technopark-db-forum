package post

import "technopark-db-forum/internal/model"

type Repository interface {
	FindById(id string) (*model.PostFull, error)
	Update(id string, message string) (*model.Post, error)
}
