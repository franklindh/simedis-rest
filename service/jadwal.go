package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrJadwalConflict      = errors.New("jadwal slot for this doctor at this time already exists")
	ErrWaktuSelesaiInvalid = errors.New("waktu_selesai must be after waktu_mulai")
)

type JadwalRepository interface {
	Create(jadwal model.Jadwal) (model.Jadwal, error)
	GetAll(params repository.ParamsGetAllJadwal) ([]model.Jadwal, pagination.Metadata, error)
	GetById(id int) (model.Jadwal, error)
	Update(id int, jadwal model.Jadwal) (model.Jadwal, error)
	Delete(id int) error
}

type JadwalService struct {
	repo JadwalRepository
}

func NewJadwalService(repo JadwalRepository) *JadwalService {
	return &JadwalService{repo: repo}
}

func (s *JadwalService) CreateJadwal(ctx context.Context, req model.JadwalRequest) (model.JadwalResponse, error) {

	layout := "15:04"
	startTime, err := time.Parse(layout, req.WaktuMulai)
	if err != nil {
		return model.JadwalResponse{}, errors.New("invalid format for waktu_mulai")
	}
	endTime, err := time.Parse(layout, req.WaktuSelesai)
	if err != nil {
		return model.JadwalResponse{}, errors.New("invalid format for waktu_selesai")
	}

	if !endTime.After(startTime) {
		return model.JadwalResponse{}, ErrWaktuSelesaiInvalid
	}

	jadwalInput := req.ToModel()
	createdJadwal, err := s.repo.Create(jadwalInput)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.JadwalResponse{}, ErrJadwalConflict
		}
		return model.JadwalResponse{}, fmt.Errorf("failed to create jadwal: %w", err)
	}

	fullJadwal, err := s.repo.GetById(createdJadwal.ID)
	if err != nil {
		return model.JadwalResponse{}, err
	}

	return model.ToJadwalResponse(fullJadwal), nil
}

func (s *JadwalService) GetAllJadwal(ctx context.Context, params repository.ParamsGetAllJadwal) ([]model.JadwalResponse, pagination.Metadata, error) {
	allJadwal, metadata, err := s.repo.GetAll(params)
	if err != nil {
		return nil, metadata, fmt.Errorf("failed to get all jadwal: %w", err)
	}

	return model.ToJadwalResponseList(allJadwal), metadata, nil
}

func (s *JadwalService) GetJadwalByID(ctx context.Context, id int) (model.JadwalResponse, error) {
	jadwal, err := s.repo.GetById(id)
	if err != nil {
		return model.JadwalResponse{}, err
	}
	return model.ToJadwalResponse(jadwal), nil
}

func (s *JadwalService) UpdateJadwal(ctx context.Context, id int, req model.JadwalRequest) (model.JadwalResponse, error) {

	layout := "15:04"
	startTime, err := time.Parse(layout, req.WaktuMulai)
	if err != nil {
		return model.JadwalResponse{}, errors.New("invalid format for waktu_mulai")
	}
	endTime, err := time.Parse(layout, req.WaktuSelesai)
	if err != nil {
		return model.JadwalResponse{}, errors.New("invalid format for waktu_selesai")
	}
	if !endTime.After(startTime) {
		return model.JadwalResponse{}, ErrWaktuSelesaiInvalid
	}

	jadwalUpdate := req.ToModel()
	updatedJadwal, err := s.repo.Update(id, jadwalUpdate)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.JadwalResponse{}, ErrJadwalConflict
		}
		return model.JadwalResponse{}, err
	}

	return model.ToJadwalResponse(updatedJadwal), nil
}

func (s *JadwalService) DeleteJadwal(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
