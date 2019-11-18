package forumHttp

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type ForumHandler struct {
}

func NewForumHandler(m *mux.Router) {
	handler := &ForumHandler{}

	m.HandleFunc("/forum", handler.HandleForum).Methods(http.MethodGet)
}

func (h *ForumHandler) HandleForum(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello: %v\n", "World")
}
