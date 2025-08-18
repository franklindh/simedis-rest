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
	ErrPoliRestored = errors.New("data restored successfully")
)

type PoliService struct {
	repo *repository.PoliRepository
}

func NewPoliService(repo *repository.PoliRepository) *PoliService {
	return &PoliService{repo: repo}
}

func (s *PoliService) CreatePoli(ctx context.Context, req model.CreatePoliRequest) (model.Poli, error) {
	poliInput := model.Poli{
		Nama:   req.Name,
		Status: req.Status,
	}

	createdPoli, err := s.repo.Create(poliInput)
	if err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.Poli{}, ErrPoliConflict
		}
		return model.Poli{}, fmt.Errorf("failed to create poli: %w", err)
	}
	return createdPoli, nil
}

func (s *PoliService) GetAllPolis(ctx context.Context) ([]model.Poli, error) {
	return s.repo.GetAll()
}

func (s *PoliService) GetPoliByID(ctx context.Context, id int) (model.Poli, error) {
	return s.repo.GetByID(id)
}

func (s *PoliService) UpdatePoli(ctx context.Context, id int, req model.UpdatePoliRequest) (model.Poli, error) {

	existing, err := s.repo.FindByName(req.Name)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return model.Poli{}, fmt.Errorf("database error on check: %w", err)
	}
	if err == nil && existing.ID != id {
		return model.Poli{}, ErrPoliConflict
	}

	poliUpdate := model.Poli{
		Nama:   req.Name,
		Status: req.Status,
	}
	return s.repo.Update(id, poliUpdate)
}

func (s *PoliService) DeletePoli(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
