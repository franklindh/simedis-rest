package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrForeignKey      = errors.New("invalid jadwal_id or pasien_id")
	ErrAntrianExists   = errors.New("pasien sudah terdaftar di jadwal ini")
	ErrScheduleOverlap = errors.New("pasien memiliki jadwal lain yang tumpang tindih pada waktu yang sama")
)

type AntrianService struct {
	repo       *repository.AntrianRepository
	jadwalRepo *repository.JadwalRepository
}

func NewAntrianService(repo *repository.AntrianRepository, jadwalRepo *repository.JadwalRepository) *AntrianService {
	return &AntrianService{repo: repo, jadwalRepo: jadwalRepo}
}

func (s *AntrianService) CreateAntrian(ctx context.Context, req model.CreateAntrianRequest) (model.Antrian, error) {
	jadwal, err := s.jadwalRepo.GetById(req.JadwalID)
	if err != nil {
		return model.Antrian{}, ErrForeignKey
	}

	tanggalString := jadwal.Tanggal.Format("2006-01-02")
	waktuMulaiString := jadwal.WaktuMulai.Format("15:04")
	waktuSelesaiString := jadwal.WaktuSelesai.Format("15:04")

	err = s.repo.CheckForOverlappingAntrian(req.PasienID, tanggalString, waktuMulaiString, waktuSelesaiString)
	if err == nil {
		return model.Antrian{}, ErrScheduleOverlap
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return model.Antrian{}, fmt.Errorf("error checking for overlapping schedule: %w", err)
	}

	err = s.repo.CheckAntrian(req.PasienID, req.JadwalID)
	if err == nil {
		return model.Antrian{}, ErrAntrianExists
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return model.Antrian{}, fmt.Errorf("error checking existing antrian: %w", err)
	}

	jadwal, err = s.jadwalRepo.GetById(req.JadwalID)
	if err != nil {
		return model.Antrian{}, ErrForeignKey
	}

	initial := strings.ToUpper(string(jadwal.Poli.Nama[0]))
	nomorAntrian := fmt.Sprintf("%s%d", initial, (time.Now().Unix()%1000)+100)

	antrian := model.Antrian{
		JadwalID:     req.JadwalID,
		PasienID:     req.PasienID,
		Prioritas:    req.Prioritas,
		Status:       "Menunggu",
		NomorAntrian: nomorAntrian,
	}

	createdAntrian, err := s.repo.Create(antrian)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23503" {
			return model.Antrian{}, ErrForeignKey
		}
		return model.Antrian{}, fmt.Errorf("failed to create antrian: %w", err)
	}

	return createdAntrian, nil
}

func (s *AntrianService) GetAllAntrianDetails(ctx context.Context, params repository.ParamsGetAllAntrian) ([]model.AntrianDetail, pagination.Metadata, error) {
	allAntrian, metadata, err := s.repo.GetAll(params)
	if err != nil {
		return nil, pagination.Metadata{}, err
	}

	var responseData []model.AntrianDetail
	for _, antrian := range allAntrian {
		detail := model.AntrianDetail{
			ID:           antrian.ID,
			NomorAntrian: antrian.NomorAntrian,
			Prioritas:    antrian.Prioritas,
			Status:       antrian.Status,
			Jadwal: struct {
				ID      int    `json:"id"`
				Tanggal string `json:"tanggal"`
				Poli    struct {
					Name string `json:"name"`
				} `json:"poli"`
				Dokter struct {
					Name string `json:"name"`
				} `json:"dokter"`
			}{
				ID:      antrian.Jadwal.ID,
				Tanggal: antrian.Jadwal.Tanggal.Format("2006-01-02"),
				Poli: struct {
					Name string `json:"name"`
				}{Name: antrian.Jadwal.Poli.Nama},
				Dokter: struct {
					Name string `json:"name"`
				}{Name: antrian.Jadwal.Petugas.Nama},
			},
			Pasien: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{
				ID:   antrian.Pasien.ID,
				Name: antrian.Pasien.NamaPasien,
			},
		}
		responseData = append(responseData, detail)
	}

	return responseData, metadata, nil
}

func (s *AntrianService) GetAntrianByID(ctx context.Context, id int) (model.Antrian, error) {
	return s.repo.GetByID(id)
}

func (s *AntrianService) UpdateAntrian(ctx context.Context, id int, req model.UpdateAntrianRequest) (model.Antrian, error) {

	antrianUpdate := model.Antrian{
		Prioritas: req.Prioritas,
		Status:    req.Status,
	}
	return s.repo.Update(id, antrianUpdate)
}

func (s *AntrianService) DeleteAntrian(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
