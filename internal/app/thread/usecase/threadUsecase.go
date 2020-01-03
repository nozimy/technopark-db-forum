package threadUsecase

import (
	"github.com/nozimy/technopark-db-forum/internal/app/thread"
	"github.com/nozimy/technopark-db-forum/internal/app/user"
	"github.com/nozimy/technopark-db-forum/internal/app/validation"
	"github.com/nozimy/technopark-db-forum/internal/model"
	"strconv"
)

type ThreadUsecase struct {
	threadRep thread.Repository
	userRep   user.Repository
}

func (t ThreadUsecase) Vote(threadSlugOrId string, vote *model.Vote) (*model.Thread, error) {
	threadObj, err := t.FindByIdOrSlug(threadSlugOrId)
	if err != nil {
		return nil, err
	}

	userObj, err := t.userRep.FindByNickname(vote.Nickname)
	if userObj == nil || err != nil {
		return nil, err
	}

	threadObj, err = t.threadRep.Vote(threadObj, vote)
	if err != nil {
		return nil, err
	}

	return threadObj, nil
}

func (t ThreadUsecase) GetThreadPosts(threadSlugOrId string, params map[string][]string) (model.Posts, error) {
	threadObj, err := t.FindByIdOrSlug(threadSlugOrId)
	if err != nil {
		return nil, err
	}

	limits := params["limit"]
	limit := "100"
	if len(limits) >= 1 {
		limit = limits[0]
	}
	descs := params["desc"]
	desc := ""
	if len(descs) >= 1 && descs[0] == "true" {
		desc = "desc"
	}
	sinces := params["since"]
	since := ""
	if len(sinces) >= 1 {
		since = sinces[0]
	}
	sorts := params["sort"]
	sort := "flat"
	if len(sorts) >= 1 {
		sort = sorts[0]
	}

	posts, err := t.threadRep.GetThreadPosts(threadObj, limit, desc, since, sort)
	if err != nil {
		return nil, err
	}

	if sort == "tree" || sort == "parent_tree" {
		//posts = makeTree(posts)
	}

	return posts, nil
}

func (t ThreadUsecase) UpdateThread(threadSlugOrId string, threadUpdate *model.ThreadUpdate) (*model.Thread, error) {
	id, _ := strconv.Atoi(threadSlugOrId)

	threadObj, err := t.threadRep.FindByIdOrSlug(id, threadSlugOrId)
	if err != nil {
		return nil, err
	}

	if validation.IsEmptyString(threadUpdate.Title) && validation.IsEmptyString(threadUpdate.Message) {
		return threadObj, nil
	}

	if validation.IsEmptyString(threadUpdate.Title) {
		threadUpdate.Title = threadObj.Title
	}

	if validation.IsEmptyString(threadUpdate.Message) {
		threadUpdate.Message = threadObj.Message
	}

	threadObj, err = t.threadRep.UpdateThread(id, threadSlugOrId, threadUpdate)
	if err != nil {
		return nil, err
	}

	return threadObj, nil
}

func (t ThreadUsecase) FindByIdOrSlug(threadSlugOrId string) (*model.Thread, error) {
	id, _ := strconv.Atoi(threadSlugOrId)

	threadObj, err := t.threadRep.FindByIdOrSlug(id, threadSlugOrId)
	if threadObj == nil || err != nil {
		return nil, err
	}

	return threadObj, nil
}

func (t ThreadUsecase) CreatePosts(threadSlugOrId string, posts *model.Posts) (*model.Posts, int, error) {
	threadObj, err := t.FindByIdOrSlug(threadSlugOrId)
	if threadObj == nil || err != nil {
		return nil, 404, err
	}

	// todo: Maybe Goroutines?
	for _, post := range *posts {
		userObj, err := t.userRep.FindByNickname(post.Author)
		if userObj == nil || err != nil {
			return nil, 404, err
		}
	}

	posts, err = t.threadRep.CreatePosts(threadObj, posts)
	if err != nil {
		return nil, 409, err
	}

	return posts, 201, nil
}

func NewThreadUsecase(t thread.Repository, u user.Repository) thread.Usecase {
	return &ThreadUsecase{
		threadRep: t,
		userRep:   u,
	}
}

func makeTree(posts model.Posts) model.Posts {
	tree := make(model.Posts, 0)
	var parent *model.Post

	for _, p := range posts {
		if len(p.Path) == 1 {
			tree = append(tree, p)
			parent = p
		} else if len(p.Path) > 1 {
			if p.Parent == parent.ID {
				//parent.Childs = append(parent.Childs, p)
				tree = append(tree, p)
				p.ParentPointer = parent
				parent = p
			} else {
				for p.Parent != parent.ID {
					parent = parent.ParentPointer
				}
				//parent.Childs = append(parent.Childs, p)
				tree = append(tree, p)
				p.ParentPointer = parent
				parent = p
			}
		}
	}

	return tree
}
