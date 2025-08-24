package service

import (
	"context"
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

type IcdRepository interface {
	Create(icd model.Icd) (model.Icd, error)
	GetAll(params repository.ParamsGetAllIcd) ([]model.Icd, pagination.Metadata, error)
	GetByID(id int) (model.Icd, error)
	Update(id int, icd model.Icd) (model.Icd, error)
	Delete(id int) error
}

type IcdService struct {
	repo IcdRepository
}

func NewIcdService(repo IcdRepository) *IcdService {
	return &IcdService{repo: repo}
}

func (s *IcdService) CreateIcd(ctx context.Context, req model.CreateIcdRequest) (model.IcdResponse, error) {
	icd := req.ToModel()
	createdIcd, err := s.repo.Create(icd)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.IcdResponse{}, ErrIcdConflict
		}
		return model.IcdResponse{}, fmt.Errorf("failed to create ICD: %w", err)
	}
	return model.ToIcdResponse(createdIcd), nil
}

func (s *IcdService) GetAllIcd(ctx context.Context, params repository.ParamsGetAllIcd) ([]model.IcdResponse, pagination.Metadata, error) {
	icds, metadata, err := s.repo.GetAll(params)
	if err != nil {
		return nil, metadata, err
	}
	return model.ToIcdResponseList(icds), metadata, nil
}

func (s *IcdService) GetIcdByID(ctx context.Context, id int) (model.IcdResponse, error) {
	icd, err := s.repo.GetByID(id)
	if err != nil {
		return model.IcdResponse{}, err
	}
	return model.ToIcdResponse(icd), nil
}

func (s *IcdService) UpdateIcd(ctx context.Context, id int, req model.UpdateIcdRequest) (model.IcdResponse, error) {
	icdUpdate := req.ToModel()
	updatedIcd, err := s.repo.Update(id, icdUpdate)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.IcdResponse{}, ErrIcdConflict
		}
		return model.IcdResponse{}, err
	}
	return model.ToIcdResponse(updatedIcd), nil
}

func (s *IcdService) DeleteIcd(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
