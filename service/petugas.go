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
	ErrPetugasConflict     = errors.New("username already exists")
	ErrInvalidCredentials  = errors.New("invalid username or password")
	ErrOldPasswordMismatch = errors.New("old password does not match")
)

type PetugasService struct {
	repo   PetugasRepository
	config *config.Config
}

func NewPetugasService(repo PetugasRepository, cfg *config.Config) *PetugasService {
	return &PetugasService{repo: repo, config: cfg}
}

func (s *PetugasService) Login(ctx context.Context, req model.LoginPetugasRequest) (string, error) {
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("database error: %w", err)
	}

	err = utils.VerifyPassword(req.Password, user.Password)
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

func (s *PetugasService) CreatePetugas(ctx context.Context, req model.CreatePetugasRequest) (model.PetugasResponse, error) {

	petugasInput := req.ToModel()

	hashedPassword, err := utils.HashPassword(s.config.DefaultPetugasPassword)
	if err != nil {
		return model.PetugasResponse{}, fmt.Errorf("failed to hash password: %w", err)
	}

	petugasInput.Password = hashedPassword

	createdPetugas, err := s.repo.Create(petugasInput)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.PetugasResponse{}, ErrPetugasConflict
		}
		return model.PetugasResponse{}, fmt.Errorf("failed to create petugas: %w", err)
	}

	return model.ToPetugasResponse(createdPetugas), nil
}

func (s *PetugasService) GetAllPetugas(ctx context.Context, params repository.ParamsGetAllPetugas) ([]model.PetugasResponse, pagination.Metadata, error) {
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

	response := model.ToPetugasResponseList(allPetugas)
	return response, metadata, nil
}

func (s *PetugasService) GetPetugasByID(ctx context.Context, id int) (model.PetugasResponse, error) {
	petugas, err := s.repo.GetById(id)
	if err != nil {

		return model.PetugasResponse{}, err
	}
	return model.ToPetugasResponse(petugas), nil
}

func (s *PetugasService) UpdatePetugas(ctx context.Context, id int, req model.UpdatePetugasRequest) (model.PetugasResponse, error) {

	petugasUpdate := req.ToModel()

	updatedPetugas, err := s.repo.Update(id, petugasUpdate)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return model.PetugasResponse{}, ErrPetugasConflict
		}

		return model.PetugasResponse{}, err
	}

	return model.ToPetugasResponse(updatedPetugas), nil
}

func (s *PetugasService) DeletePetugas(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}

func (s *PetugasService) ChangePassword(ctx context.Context, id int, req model.ChangePasswordRequest) error {

	petugas, err := s.repo.GetById(id)
	if err != nil {
		return err
	}

	err = utils.VerifyPassword(req.OldPassword, petugas.Password)
	if err != nil {
		return ErrOldPasswordMismatch
	}

	newHashedPassword, _ := utils.HashPassword(req.NewPassword)
	return s.repo.UpdatePassword(id, newHashedPassword)
}
