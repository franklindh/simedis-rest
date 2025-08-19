package service

import (
	"context"
	"database/sql"
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
	repo *repository.JenisPemeriksaanLabRepository
}

func NewJenisPemeriksaanLabService(repo *repository.JenisPemeriksaanLabRepository) *JenisPemeriksaanLabService {
	return &JenisPemeriksaanLabService{repo: repo}
}

func (s *JenisPemeriksaanLabService) Create(ctx context.Context, req model.CreateJenisPemeriksaanLabRequest) (model.JenisPemeriksaanLab, error) {
	jenis := model.JenisPemeriksaanLab{
		NamaPemeriksaan: req.NamaPemeriksaan,
		Satuan:          sql.NullString{String: req.Satuan, Valid: req.Satuan != ""},
		NilaiRujukan:    sql.NullString{String: req.NilaiRujukan, Valid: req.NilaiRujukan != ""},
		Kriteria:        sql.NullString{String: req.Kriteria, Valid: req.Kriteria != ""},
	}
	created, err := s.repo.Create(jenis)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.JenisPemeriksaanLab{}, ErrJenisPemeriksaanConflict
		}
		return model.JenisPemeriksaanLab{}, fmt.Errorf("failed to create jenis pemeriksaan: %w", err)
	}
	return created, nil
}

func (s *JenisPemeriksaanLabService) GetAll(ctx context.Context) ([]model.JenisPemeriksaanLab, error) {
	return s.repo.GetAll()
}

func (s *JenisPemeriksaanLabService) GetByID(ctx context.Context, id int) (model.JenisPemeriksaanLab, error) {
	return s.repo.GetByID(id)
}

func (s *JenisPemeriksaanLabService) Update(ctx context.Context, id int, req model.UpdateJenisPemeriksaanLabRequest) (model.JenisPemeriksaanLab, error) {
	jenis := model.JenisPemeriksaanLab{
		NamaPemeriksaan: req.NamaPemeriksaan,
		Satuan:          sql.NullString{String: req.Satuan, Valid: req.Satuan != ""},
		NilaiRujukan:    sql.NullString{String: req.NilaiRujukan, Valid: req.NilaiRujukan != ""},
		Kriteria:        sql.NullString{String: req.Kriteria, Valid: req.Kriteria != ""},
	}
	updated, err := s.repo.Update(id, jenis)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.JenisPemeriksaanLab{}, ErrJenisPemeriksaanConflict
		}
		return model.JenisPemeriksaanLab{}, err
	}
	return updated, nil
}

func (s *JenisPemeriksaanLabService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
