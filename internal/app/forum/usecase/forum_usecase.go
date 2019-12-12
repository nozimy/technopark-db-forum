package forumUsecase

import (
	"github.com/pkg/errors"
	"log"
	"technopark-db-forum/internal/app/forum"
	"technopark-db-forum/internal/app/thread"
	"technopark-db-forum/internal/app/user"
	"technopark-db-forum/internal/model"
)

type ForumUsecase struct {
	forumRep  forum.Repository
	userRep   user.Repository
	threadRep thread.Repository
}

func (f ForumUsecase) GetThreadsByForum(forumSlug string, params map[string][]string) (model.Threads, int, error) {
	forumObj, err := f.Find(forumSlug)
	if forumObj == nil || err != nil {
		return nil, 404, err
	}

	threads, err := f.forumRep.FindForumThreads(forumSlug, params)
	if err != nil {
		return nil, 404, err
	}

	return threads, 200, nil
}

func (f ForumUsecase) GetUsersByForum(forumSlug string, params map[string][]string) (model.Users, int, error) {
	forumObj, err := f.Find(forumSlug)
	if forumObj == nil || err != nil {
		return nil, 404, err
	}

	users, err := f.forumRep.FindForumUsers(forumSlug, params)
	if err != nil {
		return nil, 404, err
	}

	return users, 200, nil
}

func (f ForumUsecase) CreateThread(forumSlug string, newThread *model.NewThread) (*model.Thread, int, error) {
	userObj, err := f.userRep.FindByNickname(newThread.Author)
	if userObj == nil || err != nil {
		return nil, 404, err
	}

	forumObj, err := f.Find(forumSlug)
	if forumObj == nil || err != nil {
		return nil, 404, err
	}

	newThread.Forum = forumSlug

	threadObj, err := f.threadRep.FindByIdOrSlug(0, newThread.Slug)
	if threadObj != nil {
		return threadObj, 409, err
	}

	threadObj, err = f.threadRep.CreateThread(newThread)
	if err != nil {
		return nil, 409, err
	}

	return threadObj, 201, nil
}

func (f ForumUsecase) CreateForum(data *model.Forum) (*model.Forum, int, error) {
	userObj, err := f.userRep.FindByNickname(data.User)
	if userObj == nil || err != nil {
		return nil, 404, err
	}

	if err := f.forumRep.Create(data); err != nil {
		forumObj, err := f.forumRep.Find(data.Slug)
		if err != nil {
			log.Println(err)
			return nil, 409, err
		}

		return forumObj, 409, err
	}

	return data, 201, nil
}

func (f ForumUsecase) Find(slug string) (*model.Forum, error) {
	forumObj, err := f.forumRep.Find(slug)

	if err != nil {
		return nil, errors.Wrap(err, "forumRep.Find()")
	}

	return forumObj, nil
}

func NewForumUsecase(f forum.Repository, u user.Repository, t thread.Repository) forum.Usecase {
	return &ForumUsecase{
		forumRep:  f,
		userRep:   u,
		threadRep: t,
	}
}
