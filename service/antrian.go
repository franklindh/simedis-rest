package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrForeignKey      = errors.New("invalid jadwal_id or pasien_id")
	ErrAntrianExists   = errors.New("pasien sudah terdaftar di jadwal ini")
	ErrScheduleOverlap = errors.New("pasien memiliki jadwal lain yang tumpang tindih")
)

type AntrianService struct {
	repo       AntrianRepository
	jadwalRepo JadwalRepository
}

func NewAntrianService(repo AntrianRepository, jadwalRepo JadwalRepository) *AntrianService {
	return &AntrianService{repo: repo, jadwalRepo: jadwalRepo}
}

func (s *AntrianService) CreateAntrian(ctx context.Context, req model.CreateAntrianRequest) (model.AntrianResponse, error) {

	jadwal, err := s.jadwalRepo.GetById(req.JadwalID)
	if err != nil {
		return model.AntrianResponse{}, ErrForeignKey
	}

	isOverlap, err := s.repo.CheckForOverlappingAntrian(req.PasienID, jadwal.Tanggal, jadwal.WaktuMulai, jadwal.WaktuSelesai)
	if err != nil {
		return model.AntrianResponse{}, fmt.Errorf("error checking for overlapping schedule: %w", err)
	}
	if isOverlap {
		return model.AntrianResponse{}, ErrScheduleOverlap
	}

	isExist, err := s.repo.CheckAntrian(req.PasienID, req.JadwalID)
	if err != nil {
		return model.AntrianResponse{}, fmt.Errorf("error checking existing antrian: %w", err)
	}
	if isExist {
		return model.AntrianResponse{}, ErrAntrianExists
	}

	count, err := s.repo.CountTodayByJadwal(req.JadwalID)
	if err != nil {
		return model.AntrianResponse{}, fmt.Errorf("error counting antrian for today's schedule: %w", err)
	}
	initial := strings.ToUpper(string(jadwal.Poli.Nama[0]))
	nomorAntrian := fmt.Sprintf("%s%d", initial, count+1)

	antrian := req.ToModel(nomorAntrian)
	createdAntrian, err := s.repo.Create(antrian)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23503" {
			return model.AntrianResponse{}, ErrForeignKey
		}
		return model.AntrianResponse{}, fmt.Errorf("failed to create antrian: %w", err)
	}

	return model.ToAntrianResponse(createdAntrian), nil
}

func (s *AntrianService) GetAllAntrian(ctx context.Context, params repository.ParamsGetAllAntrian) ([]model.AntrianResponse, pagination.Metadata, error) {
	allAntrian, metadata, err := s.repo.GetAll(params)
	if err != nil {
		return nil, metadata, err
	}
	return model.ToAntrianResponseList(allAntrian), metadata, nil
}

func (s *AntrianService) GetAntrianByID(ctx context.Context, id int) (model.AntrianResponse, error) {
	antrian, err := s.repo.GetByID(id)
	if err != nil {
		return model.AntrianResponse{}, err
	}
	return model.ToAntrianResponse(antrian), nil
}

func (s *AntrianService) UpdateAntrian(ctx context.Context, id int, req model.UpdateAntrianRequest) (model.AntrianResponse, error) {
	antrianUpdate := req.ToModel()
	updatedAntrian, err := s.repo.Update(id, antrianUpdate)
	if err != nil {
		return model.AntrianResponse{}, err
	}
	return model.ToAntrianResponse(updatedAntrian), nil
}

func (s *AntrianService) DeleteAntrian(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
