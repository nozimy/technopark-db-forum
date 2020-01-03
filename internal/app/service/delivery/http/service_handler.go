package serviceHttp

import (
	"github.com/gorilla/mux"
	"net/http"
	"technopark-db-forum/internal/app/respond"
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
	err := h.ServiceUsecase.ClearAll()

	if err != nil {
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, nil)
}

func (h *ServiceHandler) HandleServiceGetStatus(w http.ResponseWriter, r *http.Request) {
	status, err := h.ServiceUsecase.GetStatus()

	if err != nil || status == nil {
		respond.Error(w, r, http.StatusNotFound, err)
		return
	}

	respond.Respond(w, r, http.StatusOK, status)
}
