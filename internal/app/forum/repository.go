package forum

import "technopark-db-forum/internal/model"

type Repository interface {
	Create(forum *model.Forum) error
	Find(int64) (*model.Forum, error)
}
