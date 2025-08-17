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

type JadwalService struct {
	repo *repository.JadwalRepository
}

func NewJadwalService(repo *repository.JadwalRepository) *JadwalService {
	return &JadwalService{repo: repo}
}

func (s *JadwalService) CreateJadwal(ctx context.Context, req model.JadwalRequest) (model.Jadwal, error) {
	layout := "15:04"
	startTime, err := time.Parse(layout, req.WaktuMulai)
	if err != nil {

		return model.Jadwal{}, errors.New("invalid format for waktu_mulai")
	}

	endTime, err := time.Parse(layout, req.WaktuSelesai)
	if err != nil {
		return model.Jadwal{}, errors.New("invalid format for waktu_selesai")
	}

	if !endTime.After(startTime) {
		return model.Jadwal{}, ErrWaktuSelesaiInvalid
	}

	jadwalInput := req.ToModel()

	createdJadwal, err := s.repo.Create(jadwalInput)
	if err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.Jadwal{}, ErrJadwalConflict
		}
		return model.Jadwal{}, fmt.Errorf("failed to create jadwal: %w", err)
	}
	return s.repo.GetByID(createdJadwal.ID)
}

func (s *JadwalService) GetAllJadwalDetails(ctx context.Context, params repository.ParamsGetAllJadwal) ([]model.JadwalDetail, pagination.Metadata, error) {

	allJadwal, metadata, err := s.repo.GetAll(params)
	if err != nil {
		return nil, pagination.Metadata{}, fmt.Errorf("failed to get all jadwal: %w", err)
	}

	var responseData []model.JadwalDetail
	for _, jadwal := range allJadwal {
		detail := model.JadwalDetail{
			ID:           jadwal.ID,
			Tanggal:      jadwal.Tanggal,
			WaktuMulai:   jadwal.WaktuMulai,
			WaktuSelesai: jadwal.WaktuSelesai,
			Keterangan:   jadwal.Keterangan,
			Petugas: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{
				ID:   jadwal.Petugas.ID,
				Name: jadwal.Petugas.Nama,
			},
			Poli: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{
				ID:   jadwal.Poli.ID,
				Name: jadwal.Poli.Nama,
			},
		}
		responseData = append(responseData, detail)
	}

	return responseData, metadata, nil
}

func (s *JadwalService) GetJadwalByID(ctx context.Context, id int) (model.Jadwal, error) {
	return s.repo.GetByID(id)
}

func (s *JadwalService) UpdateJadwal(ctx context.Context, id int, req model.JadwalRequest) (model.Jadwal, error) {

	jadwalUpdate := req.ToModel()
	updatedJadwal, err := s.repo.Update(id, jadwalUpdate)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.Jadwal{}, ErrJadwalConflict
		}
		return model.Jadwal{}, err
	}
	return updatedJadwal, nil
}

func (s *JadwalService) DeleteJadwal(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
