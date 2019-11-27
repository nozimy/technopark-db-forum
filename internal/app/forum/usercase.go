package forum

import "technopark-db-forum/internal/model"

type Usecase interface {
	CreateForum(*model.Forum) error
	Find(int64) (*model.Forum, error)
}
