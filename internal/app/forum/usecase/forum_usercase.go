package forumUsecase

import (
	"github.com/pkg/errors"
	"technopark-db-forum/internal/app/forum"
	"technopark-db-forum/internal/model"
)

type ForumUsecase struct {
	forumRep forum.Repository
}

func (f ForumUsecase) CreateForum(data *model.Forum) error {
	if err := f.forumRep.Create(data); err != nil {
		return errors.Wrap(err, "CreateForum<-forumRep.Create()")
	}

	return nil
}

func (f ForumUsecase) Find(id int64) (*model.Forum, error) {
	forumObj, err := f.forumRep.Find(id)

	if err != nil {
		return nil, errors.Wrap(err, "forumRep.Find()")
	}

	return forumObj, nil
}

func NewForumUsecase(f forum.Repository) forum.Usecase {
	return &ForumUsecase{
		forumRep: f,
	}
}
