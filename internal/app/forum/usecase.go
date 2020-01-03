package forum

import "github.com/nozimy/technopark-db-forum/internal/model"

type Usecase interface {
	CreateForum(*model.Forum) (*model.Forum, int, error)
	Find(slug string) (*model.Forum, error)
	CreateThread(string, *model.NewThread) (*model.Thread, int, error)
	GetUsersByForum(forumSlug string, params map[string][]string) (model.Users, int, error)
	GetThreadsByForum(forumSlug string, params map[string][]string) (model.Threads, int, error)
}
