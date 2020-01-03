package service

import "technopark-db-forum/internal/model"

type Usecase interface {
	GetStatus() (*model.Status, error)
	ClearAll() error
}
