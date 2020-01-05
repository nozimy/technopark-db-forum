package forum

import (
	"github.com/nozimy/technopark-db-forum/internal/model"
)

type Repository interface {
	Create(forum *model.Forum) error
	Find(slug string) (*model.Forum, error)
	FindForumUsers(forumObj *model.Forum, params map[string][]string) (model.Users, error)
	FindForumThreads(forumSlug string, params map[string][]string) (model.Threads, error)
}
