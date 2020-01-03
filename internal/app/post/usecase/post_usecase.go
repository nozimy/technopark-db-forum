package postUsecase

import (
	"github.com/nozimy/technopark-db-forum/internal/app/post"
	"github.com/nozimy/technopark-db-forum/internal/app/validation"
	"github.com/nozimy/technopark-db-forum/internal/model"
	"github.com/pkg/errors"
	"strings"
)

type PostUsecase struct {
	postRep post.Repository
}

func (p PostUsecase) Update(id string, message string) (*model.Post, error) {
	postFullObj, err := p.postRep.FindById(id, false, false, false)
	if err != nil {
		return nil, errors.Wrap(err, "postRep.FindById()")
	}

	if validation.IsEmptyString(message) || postFullObj.Post.Message == message {
		return postFullObj.Post, nil
	}

	postObj, err := p.postRep.Update(id, message)

	if err != nil {
		return nil, errors.Wrap(err, "postRep.Update()")
	}

	return postObj, nil
}

func (p PostUsecase) FindById(id string, params map[string][]string) (*model.PostFull, error) {
	related := params["related"]
	includeUser := false
	includeForum := false
	includeThread := false
	if len(related) >= 1 {
		splitRelated := strings.Split(related[0], ",")

		if contains(splitRelated, "user") {
			includeUser = true
		}
		if contains(splitRelated, "forum") {
			includeForum = true
		}
		if contains(splitRelated, "thread") {
			includeThread = true
		}
	}

	postObj, err := p.postRep.FindById(id, includeUser, includeForum, includeThread)

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

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}
