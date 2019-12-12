package serviceHttp

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"technopark-db-forum/internal/app/service"
)

type ServiceHandler struct {
	ServiceUsecase service.Usecase
}

func NewServiceHandler(m *mux.Router, u service.Usecase) {
	handler := &ServiceHandler{
		ServiceUsecase: u,
	}

	m.HandleFunc("/service/clear", handler.HandleServiceClear).Methods(http.MethodPost)
	m.HandleFunc("/service/status", handler.HandleServiceGetStatus).Methods(http.MethodGet)
}

func (h *ServiceHandler) HandleServiceClear(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello: %v\n", "World")
}

func (h *ServiceHandler) HandleServiceGetStatus(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello: %v\n", "World")
}
