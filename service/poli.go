package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
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

func (s *PoliService) CreateOrRestorePoli(ctx context.Context, poliInput model.Poli) (model.Poli, error) {

	existingPoli, err := s.repo.FindByNameIncludingDeleted(poliInput.Nama)
	if err != nil && err != sql.ErrNoRows {

		return model.Poli{}, fmt.Errorf("database error: %w", err)
	}

	if err == nil {

		if existingPoli.DeletedAt.Valid {
			restoredPoli, restoreErr := s.repo.Restore(existingPoli.ID)
			if restoreErr != nil {
				return model.Poli{}, fmt.Errorf("failed to restore poli: %w", restoreErr)
			}

			return restoredPoli, ErrPoliRestored
		}

		return model.Poli{}, ErrPoliConflict
	}

	createdPoli, createErr := s.repo.Create(poliInput)
	if createErr != nil {
		return model.Poli{}, fmt.Errorf("failed to create poli: %w", createErr)
	}
	return createdPoli, nil
}

func (s *PoliService) GetAllPolis(ctx context.Context) ([]model.Poli, error) {
	return s.repo.GetAll()
}

func (s *PoliService) GetPoliByID(ctx context.Context, id int) (model.Poli, error) {
	return s.repo.GetByID(id)
}

func (s *PoliService) UpdatePoli(ctx context.Context, id int, poliInput model.Poli) (model.Poli, error) {

	existing, err := s.repo.FindByName(poliInput.Nama)
	if err != nil && err != repository.ErrNotFound {
		return model.Poli{}, fmt.Errorf("database error on check: %w", err)
	}

	if err == nil && existing.ID != id {
		return model.Poli{}, ErrPoliConflict
	}

	return s.repo.Update(id, poliInput)
}

func (s *PoliService) DeletePoli(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
