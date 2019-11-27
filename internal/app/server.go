package apiserver

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	forumHttp "technopark-db-forum/internal/app/forum/delivery/http"
	"technopark-db-forum/internal/app/forum/repository"
	forumUsecase "technopark-db-forum/internal/app/forum/usecase"
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
		Mux:    mux.NewRouter(),
		Config: config,
	}

	return server, nil
}

func (s *Server) ConfigureServer(db *sql.DB) {
	forumRep := forumRepository.NewForumRepository(db)

	forumUse := forumUsecase.NewForumUsecase(forumRep)

	forumHttp.NewForumHandler(s.Mux, forumUse)
}
