package postHttp

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"technopark-db-forum/internal/app/post"
)

type PostHandler struct {
	PostUsecase post.Usecase
}

func NewPostHandler(m *mux.Router, u post.Usecase) {
	handler := &PostHandler{
		PostUsecase: u,
	}

	m.HandleFunc("/post/{id}/details", handler.HandleUpdatePost).Methods(http.MethodPost, http.MethodOptions)
	m.HandleFunc("/post/{id}/details", handler.HandleGetPostDetails).Methods(http.MethodGet)
}

func (h *PostHandler) HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello: %v\n", "World")
}

func (h *PostHandler) HandleGetPostDetails(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello: %v\n", "World")
}
