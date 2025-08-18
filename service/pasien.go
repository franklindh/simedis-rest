package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrPasienConflict = errors.New("data with the same NIK, username, or nomor kartu jaminan already exists")
)

type PasienService struct {
	repo *repository.PasienRepository
}

func NewPasienService(repo *repository.PasienRepository) *PasienService {
	return &PasienService{repo: repo}
}

func (s *PasienService) CreatePasien(ctx context.Context, req model.CreatePasienRequest) (model.Pasien, error) {
	var username, password string

	if req.UsernamePasien == "" {
		username = req.NIK
	} else {
		username = req.UsernamePasien
	}

	if req.Password == "" {
		password = req.NIK
	} else {
		password = req.Password
	}
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return model.Pasien{}, fmt.Errorf("failed to hash password: %w", err)
	}

	lastID, _ := s.repo.GetLastID()
	noRekamMedis := fmt.Sprintf("RM-%s-%04d", time.Now().Format("20060102"), lastID+1)

	pasien := req.ToModel(username, hashedPassword, noRekamMedis)

	createdPasien, err := s.repo.Create(pasien)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.Pasien{}, ErrPasienConflict
		}
		return model.Pasien{}, fmt.Errorf("failed to create patient: %w", err)
	}

	return createdPasien, nil
}

func (s *PasienService) GetAllPasien(ctx context.Context, params repository.ParamsGetAllPasien) ([]model.Pasien, pagination.Metadata, error) {
	return s.repo.GetAll(params)
}

func (s *PasienService) GetPasienByID(ctx context.Context, id int) (model.Pasien, error) {
	return s.repo.GetByID(id)
}

func (s *PasienService) UpdatePasien(ctx context.Context, id int, req model.UpdatePasienRequest) (model.Pasien, error) {
	pasienUpdate := req.ToModel()

	updatedPasien, err := s.repo.Update(id, pasienUpdate)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.Pasien{}, ErrPasienConflict
		}
		return model.Pasien{}, err
	}
	return updatedPasien, nil
}

func (s *PasienService) DeletePasien(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
