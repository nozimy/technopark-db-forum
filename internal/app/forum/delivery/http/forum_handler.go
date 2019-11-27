package forumHttp

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"technopark-db-forum/internal/app/forum"
	"technopark-db-forum/internal/model"
)

type ForumHandler struct {
	ForumUsecase forum.Usecase
}

func NewForumHandler(m *mux.Router, fu forum.Usecase) {
	handler := &ForumHandler{
		ForumUsecase: fu,
	}

	m.HandleFunc("/forum", handler.HandleForum).Methods(http.MethodGet)
	m.HandleFunc("/forum-create", handler.HandleCreateForum).Methods(http.MethodGet)
}

func (h *ForumHandler) HandleCreateForum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	f := &model.Forum{
		Name: "",
	}
	if err := h.ForumUsecase.CreateForum(f); err != nil {
		err = errors.Wrapf(err, "HandleCreateForum<-CreateForum: ")
		return
	}

	_, _ = fmt.Fprintf(w, "Hello: %v\n", "World")
}

func (h *ForumHandler) HandleForum(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello: %v\n", "World")
}
