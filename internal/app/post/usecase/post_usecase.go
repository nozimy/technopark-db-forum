package postUsecase

import (
	"github.com/pkg/errors"
	"technopark-db-forum/internal/app/post"
	"technopark-db-forum/internal/model"
)

type PostUsecase struct {
	postRep post.Repository
}

func (p PostUsecase) Update(id string, message string) (*model.Post, error) {
	postObj, err := p.postRep.Update(id, message)

	if err != nil {
		return nil, errors.Wrap(err, "postRep.Update()")
	}

	return postObj, nil
}

func (p PostUsecase) FindById(id string) (*model.PostFull, error) {
	postObj, err := p.postRep.FindById(id)

	if err != nil {
		return nil, errors.Wrap(err, "postRep.FindById()")
	}

	return postObj, nil
}

func NewPostUsecase(p post.Repository) post.Usecase {
	return &PostUsecase{
		postRep: p,
	}
}
