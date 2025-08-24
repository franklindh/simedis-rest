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
	ErrJenisPemeriksaanConflict = errors.New("jenis pemeriksaan with that name already exists")
)

type JenisPemeriksaanLabService struct {
	repo JenisPemeriksaanLabRepository
}

func NewJenisPemeriksaanLabService(repo JenisPemeriksaanLabRepository) *JenisPemeriksaanLabService {
	return &JenisPemeriksaanLabService{repo: repo}
}

func (s *JenisPemeriksaanLabService) Create(ctx context.Context, req model.CreateJenisPemeriksaanLabRequest) (model.JenisPemeriksaanLabResponse, error) {
	jenis := req.ToModel()

	created, err := s.repo.Create(jenis)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.JenisPemeriksaanLabResponse{}, ErrJenisPemeriksaanConflict
		}
		return model.JenisPemeriksaanLabResponse{}, fmt.Errorf("failed to create jenis pemeriksaan: %w", err)
	}

	return model.ToJenisPemeriksaanLabResponse(created), nil
}

func (s *JenisPemeriksaanLabService) GetAll(ctx context.Context) ([]model.JenisPemeriksaanLabResponse, error) {
	list, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return model.ToJenisPemeriksaanLabResponseList(list), nil
}

func (s *JenisPemeriksaanLabService) GetByID(ctx context.Context, id int) (model.JenisPemeriksaanLabResponse, error) {
	jenis, err := s.repo.GetById(id)
	if err != nil {
		return model.JenisPemeriksaanLabResponse{}, err
	}
	return model.ToJenisPemeriksaanLabResponse(jenis), nil
}

func (s *JenisPemeriksaanLabService) Update(ctx context.Context, id int, req model.UpdateJenisPemeriksaanLabRequest) (model.JenisPemeriksaanLabResponse, error) {

	existing, err := s.repo.FindByName(req.NamaPemeriksaan)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return model.JenisPemeriksaanLabResponse{}, fmt.Errorf("database error on check: %w", err)
	}
	if err == nil && existing.ID != id {
		return model.JenisPemeriksaanLabResponse{}, ErrJenisPemeriksaanConflict
	}

	jenis := req.ToModel()
	updated, err := s.repo.Update(id, jenis)
	if err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.JenisPemeriksaanLabResponse{}, ErrJenisPemeriksaanConflict
		}
		return model.JenisPemeriksaanLabResponse{}, err
	}

	return model.ToJenisPemeriksaanLabResponse(updated), nil
}

func (s *JenisPemeriksaanLabService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
