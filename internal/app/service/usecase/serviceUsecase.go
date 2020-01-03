package serviceUsecase

import (
	"github.com/nozimy/technopark-db-forum/internal/app/service"
	"github.com/nozimy/technopark-db-forum/internal/model"
	"github.com/pkg/errors"
)

type ServiceUsecase struct {
	rep service.Repository
}

func (s ServiceUsecase) ClearAll() error {
	err := s.rep.ClearAll()

	if err != nil {
		return errors.Wrap(err, "ClearAll()")
	}

	return nil
}

func (s ServiceUsecase) GetStatus() (*model.Status, error) {
	status, err := s.rep.GetStatus()

	if err != nil {
		return nil, errors.Wrap(err, "GetStatus()")
	}

	return status, nil
}

func NewServiceUsecase(rep service.Repository) service.Usecase {
	return &ServiceUsecase{
		rep: rep,
	}
}
