package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrPoliConflict = errors.New("data with that name already exists")
)

type PoliRepository interface {
	Create(poli model.Poli) (model.Poli, error)
	GetAll() ([]model.Poli, error)
	GetById(id int) (model.Poli, error)
	Update(id int, poli model.Poli) (model.Poli, error)
	Delete(id int) error
	FindByName(name string) (model.Poli, error)
}

type PoliService struct {
	repo PoliRepository
}

func NewPoliService(repo PoliRepository) *PoliService {
	return &PoliService{repo: repo}
}

func (s *PoliService) CreatePoli(ctx context.Context, req model.CreatePoliRequest) (model.PoliResponse, error) {

	poliInput := req.ToModel()

	createdPoli, err := s.repo.Create(poliInput)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.PoliResponse{}, ErrPoliConflict
		}
		return model.PoliResponse{}, fmt.Errorf("failed to create poli: %w", err)
	}

	return model.ToPoliResponse(createdPoli), nil
}

func (s *PoliService) GetAllPolis(ctx context.Context) ([]model.PoliResponse, error) {
	polis, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all polis: %w", err)
	}

	return model.ToPoliResponseList(polis), nil
}

func (s *PoliService) GetPoliByID(ctx context.Context, id int) (model.PoliResponse, error) {
	poli, err := s.repo.GetById(id)
	if err != nil {

		return model.PoliResponse{}, err
	}
	return model.ToPoliResponse(poli), nil
}

func (s *PoliService) UpdatePoli(ctx context.Context, id int, req model.UpdatePoliRequest) (model.PoliResponse, error) {
	existing, err := s.repo.FindByName(req.Name)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return model.PoliResponse{}, fmt.Errorf("database error on check: %w", err)
	}
	if err == nil && existing.ID != id {
		return model.PoliResponse{}, ErrPoliConflict
	}

	poliUpdate := req.ToModel()

	updatedPoli, err := s.repo.Update(id, poliUpdate)
	if err != nil {
		return model.PoliResponse{}, err
	}

	return model.ToPoliResponse(updatedPoli), nil
}

func (s *PoliService) DeletePoli(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
