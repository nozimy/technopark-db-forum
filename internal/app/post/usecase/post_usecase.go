package postUsecase

import (
	"technopark-db-forum/internal/app/post"
)

type PostUsecase struct {
	postRep post.Repository
}

func NewPostUsecase(p post.Repository) post.Usecase {
	return &PostUsecase{
		postRep: p,
	}
}
