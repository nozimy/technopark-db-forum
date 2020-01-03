package apiserver

import (
	"database/sql"
	"github.com/gorilla/mux"
	forumHttp "github.com/nozimy/technopark-db-forum/internal/app/forum/delivery/http"
	"github.com/nozimy/technopark-db-forum/internal/app/forum/repository"
	forumUsecase "github.com/nozimy/technopark-db-forum/internal/app/forum/usecase"
	postHttp "github.com/nozimy/technopark-db-forum/internal/app/post/delivery/http"
	postRepository "github.com/nozimy/technopark-db-forum/internal/app/post/repository"
	postUsecase "github.com/nozimy/technopark-db-forum/internal/app/post/usecase"
	serviceHttp "github.com/nozimy/technopark-db-forum/internal/app/service/delivery/http"
	serviceRepository "github.com/nozimy/technopark-db-forum/internal/app/service/repository"
	serviceUsecase "github.com/nozimy/technopark-db-forum/internal/app/service/usecase"
	threadHttp "github.com/nozimy/technopark-db-forum/internal/app/thread/delivery/http"
	threadRepository "github.com/nozimy/technopark-db-forum/internal/app/thread/repository"
	threadUsecase "github.com/nozimy/technopark-db-forum/internal/app/thread/usecase"
	userHttp "github.com/nozimy/technopark-db-forum/internal/app/user/delivery/http"
	userRepository "github.com/nozimy/technopark-db-forum/internal/app/user/repository"
	userUsecase "github.com/nozimy/technopark-db-forum/internal/app/user/usecase"
	"github.com/nozimy/technopark-db-forum/internal/middleware"
	"net/http"
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
	serviceRep := serviceRepository.NewServiceRepository(db)

	forumUse := forumUsecase.NewForumUsecase(forumRep, userRep, threadRep)
	userUse := userUsecase.NewForumUsecase(userRep)
	threadUse := threadUsecase.NewThreadUsecase(threadRep, userRep)
	postUse := postUsecase.NewPostUsecase(postRep)
	serviceUse := serviceUsecase.NewServiceUsecase(serviceRep)

	s.Mux.Use(middleware.CORSMiddleware)

	forumHttp.NewForumHandler(s.Mux, forumUse)
	userHttp.NewUserHandler(s.Mux, userUse)
	threadHttp.NewThreadHandler(s.Mux, threadUse)
	postHttp.NewPostHandler(s.Mux, postUse)
	serviceHttp.NewServiceHandler(s.Mux, serviceUse)
}
