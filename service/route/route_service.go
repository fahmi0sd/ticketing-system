package route

import (
	"errors"
	"log/slog"
)

type service struct {
	logger *slog.Logger
	repo   Repository
}

type Service interface {
	GetAll(filter SearchFilter) ([]Route, error)
	GetByID(id int) (Route, error)
}

func NewService(logger *slog.Logger, repo Repository) Service {
	return &service{logger: logger, repo: repo}
}

func (s *service) GetAll(filter SearchFilter) ([]Route, error) {
	routes, err := s.repo.GetAll(filter)
	if err != nil {
		s.logger.Error("get all routes error", slog.Any("err", err))
		return nil, errors.New("failed to fetch routes")
	}

	for i := range routes {
		routes[i].Available = routes[i].Quota - routes[i].Sold
	}

	return routes, nil
}

func (s *service) GetByID(id int) (Route, error) {
	r, err := s.repo.GetByID(id)
	if err != nil {
		return Route{}, errors.New("route not found")
	}

	r.Available = r.Quota - r.Sold
	return r, nil
}
