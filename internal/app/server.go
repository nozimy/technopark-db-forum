package apiserver

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	forumHttp "technopark-db-forum/internal/app/forum/delivery/http"
	"technopark-db-forum/internal/app/forum/repository"
	forumUsecase "technopark-db-forum/internal/app/forum/usecase"
	postHttp "technopark-db-forum/internal/app/post/delivery/http"
	postRepository "technopark-db-forum/internal/app/post/repository"
	postUsecase "technopark-db-forum/internal/app/post/usecase"
	threadHttp "technopark-db-forum/internal/app/thread/delivery/http"
	threadRepository "technopark-db-forum/internal/app/thread/repository"
	threadUsecase "technopark-db-forum/internal/app/thread/usecase"
	userHttp "technopark-db-forum/internal/app/user/delivery/http"
	userRepository "technopark-db-forum/internal/app/user/repository"
	userUsecase "technopark-db-forum/internal/app/user/usecase"
	"technopark-db-forum/internal/middleware"
)

type Server struct {
	Mux    *mux.Router
	Config *Config
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

func NewServer(config *Config) (*Server, error) {
	server := &Server{
		Mux:    mux.NewRouter().PathPrefix("/api").Subrouter(),
		Config: config,
	}

	return server, nil
}

func (s *Server) ConfigureServer(db *sql.DB) {
	forumRep := forumRepository.NewForumRepository(db)
	userRep := userRepository.NewUserRepository(db)
	threadRep := threadRepository.NewThreadRepository(db)
	postRep := postRepository.NewPostRepository(db)

	forumUse := forumUsecase.NewForumUsecase(forumRep, userRep, threadRep)
	userUse := userUsecase.NewForumUsecase(userRep)
	threadUse := threadUsecase.NewThreadUsecase(threadRep)
	postUse := postUsecase.NewPostUsecase(postRep)

	s.Mux.Use(middleware.CORSMiddleware)

	forumHttp.NewForumHandler(s.Mux, forumUse)
	userHttp.NewUserHandler(s.Mux, userUse)
	threadHttp.NewThreadHandler(s.Mux, threadUse)
	postHttp.NewPostHandler(s.Mux, postUse)
}
