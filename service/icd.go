package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrIcdConflict = errors.New("ICD code already exists")
)

type IcdService struct {
	repo *repository.IcdRepository
}

func NewIcdService(repo *repository.IcdRepository) *IcdService {
	return &IcdService{repo: repo}
}

func (s *IcdService) CreateIcd(ctx context.Context, req model.CreateIcdRequest) (model.Icd, error) {
	icd := model.Icd{
		KodeIcd:      req.KodeIcd,
		NamaPenyakit: req.NamaPenyakit,
		Deskripsi:    sql.NullString{String: req.Deskripsi, Valid: req.Deskripsi != ""},
		Status:       req.Status,
	}

	createdIcd, err := s.repo.Create(icd)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.Icd{}, ErrIcdConflict
		}
		return model.Icd{}, fmt.Errorf("failed to create ICD: %w", err)
	}
	return createdIcd, nil
}

func (s *IcdService) GetAllIcd(ctx context.Context, params repository.ParamsGetAllIcd) ([]model.Icd, pagination.Metadata, error) {
	return s.repo.GetAll(params)
}

func (s *IcdService) GetIcdByID(ctx context.Context, id int) (model.Icd, error) {
	return s.repo.GetByID(id)
}

func (s *IcdService) UpdateIcd(ctx context.Context, id int, req model.UpdateIcdRequest) (model.Icd, error) {
	icdUpdate := model.Icd{
		KodeIcd:      req.KodeIcd,
		NamaPenyakit: req.NamaPenyakit,
		Deskripsi:    sql.NullString{String: req.Deskripsi, Valid: req.Deskripsi != ""},
		Status:       req.Status,
	}

	updatedIcd, err := s.repo.Update(id, icdUpdate)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.Icd{}, ErrIcdConflict
		}
		return model.Icd{}, err
	}
	return updatedIcd, nil
}

func (s *IcdService) DeleteIcd(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
