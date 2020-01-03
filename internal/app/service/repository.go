package service

import "github.com/nozimy/technopark-db-forum/internal/model"

type Repository interface {
	GetStatus() (*model.Status, error)
	ClearAll() error
}
