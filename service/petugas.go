package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/franklindh/simedis-api/internal/config"
	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrPetugasConflict    = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type PetugasService struct {
	repo   *repository.PetugasRepository
	config *config.Config
}

func NewPetugasService(repo *repository.PetugasRepository, cfg *config.Config) *PetugasService {
	return &PetugasService{repo: repo, config: cfg}
}

func (s *PetugasService) Login(ctx context.Context, username, password string) (string, error) {

	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("database error: %w", err)
	}

	err = utils.VerifyPassword(password, user.Password)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	jwtSecret := []byte(s.config.JWTSecret)
	token, err := utils.SignToken(strconv.Itoa(user.ID), user.Username, user.Role, jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *PetugasService) CreatePetugas(ctx context.Context, petugasInput model.Petugas) (model.Petugas, error) {

	defaultPassword := s.config.DefaultPetugasPassword
	encodedHash, err := utils.HashPassword(defaultPassword)
	if err != nil {
		return model.Petugas{}, fmt.Errorf("failed to hash password: %w", err)
	}
	petugasInput.Password = encodedHash

	createdPetugas, err := s.repo.Create(petugasInput)
	if err != nil {

		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.Petugas{}, ErrPetugasConflict
		}
		return model.Petugas{}, fmt.Errorf("failed to create petugas: %w", err)
	}

	createdPetugas.Password = ""
	return createdPetugas, nil
}

func (s *PetugasService) GetAllPetugas(ctx context.Context, params repository.ParamsGetAllPetugas) ([]model.Petugas, pagination.Metadata, error) {

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 5
	}

	allPetugas, metadata, err := s.repo.GetAll(params)
	if err != nil {
		return nil, pagination.Metadata{}, fmt.Errorf("failed to get all petugas: %w", err)
	}
	return allPetugas, metadata, nil
}

func (s *PetugasService) GetPetugasByID(ctx context.Context, id int) (model.Petugas, error) {
	return s.repo.GetByID(id)
}

func (s *PetugasService) UpdatePetugas(ctx context.Context, id int, petugasInput model.Petugas) (model.Petugas, error) {

	petugasInput.Password = ""
	updatedPetugas, err := s.repo.Update(id, petugasInput)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.Petugas{}, ErrPetugasConflict
		}
		return model.Petugas{}, err
	}
	return updatedPetugas, nil
}

func (s *PetugasService) DeletePetugas(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
