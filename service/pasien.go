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

type PasienRepository interface {
	Create(pasien model.Pasien) (model.Pasien, error)
	GetAll(params repository.ParamsGetAllPasien) ([]model.Pasien, pagination.Metadata, error)
	GetById(id int) (model.Pasien, error)
	Update(id int, pasien model.Pasien) (model.Pasien, error)
	Delete(id int) error
	GetLastID() (int, error)
}

type PasienService struct {
	repo PasienRepository
}

func NewPasienService(repo PasienRepository) *PasienService {
	return &PasienService{repo: repo}
}

func (s *PasienService) CreatePasien(ctx context.Context, req model.CreatePasienRequest) (model.PasienResponse, error) {

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
		return model.PasienResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	lastID, _ := s.repo.GetLastID()
	noRekamMedis := fmt.Sprintf("RM-%s-%04d", time.Now().Format("20060102"), lastID+1)

	pasien := req.ToModel(username, hashedPassword, noRekamMedis)

	createdPasien, err := s.repo.Create(pasien)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.PasienResponse{}, ErrPasienConflict
		}
		return model.PasienResponse{}, fmt.Errorf("failed to create patient: %w", err)
	}

	return model.ToPasienResponse(createdPasien), nil
}

func (s *PasienService) GetAllPasien(ctx context.Context, params repository.ParamsGetAllPasien) ([]model.PasienResponse, pagination.Metadata, error) {
	allPasien, metadata, err := s.repo.GetAll(params)
	if err != nil {
		return nil, metadata, fmt.Errorf("failed to get all pasien: %w", err)
	}

	return model.ToPasienResponseList(allPasien), metadata, nil
}

func (s *PasienService) GetPasienByID(ctx context.Context, id int) (model.PasienResponse, error) {
	pasien, err := s.repo.GetById(id)
	if err != nil {
		return model.PasienResponse{}, err
	}
	return model.ToPasienResponse(pasien), nil
}

func (s *PasienService) UpdatePasien(ctx context.Context, id int, req model.UpdatePasienRequest) (model.PasienResponse, error) {
	pasienUpdate := req.ToModel()

	updatedPasien, err := s.repo.Update(id, pasienUpdate)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.PasienResponse{}, ErrPasienConflict
		}
		return model.PasienResponse{}, err
	}

	return model.ToPasienResponse(updatedPasien), nil
}

func (s *PasienService) DeletePasien(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
