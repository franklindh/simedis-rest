package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
)

var (
	ErrPemeriksaanExists = errors.New("pemeriksaan for this antrian already exists")
)

type PemeriksaanRepository interface {
	CheckExistingPemeriksaan(antrianID int) error
	Create(pemeriksaan model.Pemeriksaan) (model.Pemeriksaan, error)
	GetById(id int) (model.Pemeriksaan, error)
	GetAllByPasienID(pasienID int) ([]model.Pemeriksaan, error)
	Update(id int, pemeriksaan model.Pemeriksaan) (model.Pemeriksaan, error)
	Delete(id int) error
}

type AntrianRepository interface {
	GetByID(id int) (model.Antrian, error)
	Update(id int, antrian model.Antrian) (model.Antrian, error)
}

type PemeriksaanService struct {
	repo        PemeriksaanRepository
	antrianRepo AntrianRepository
}

func NewPemeriksaanService(repo PemeriksaanRepository, antrianRepo AntrianRepository) *PemeriksaanService {
	return &PemeriksaanService{repo: repo, antrianRepo: antrianRepo}
}

func (s *PemeriksaanService) CreatePemeriksaan(ctx context.Context, req model.CreatePemeriksaanRequest) (model.PemeriksaanResponse, error) {

	err := s.repo.CheckExistingPemeriksaan(req.AntrianID)
	if err == nil {
		return model.PemeriksaanResponse{}, ErrPemeriksaanExists
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return model.PemeriksaanResponse{}, fmt.Errorf("error checking existing pemeriksaan: %w", err)
	}

	pemeriksaan := req.ToModel()

	createdPemeriksaan, err := s.repo.Create(pemeriksaan)
	if err != nil {
		return model.PemeriksaanResponse{}, err
	}

	antrian, err := s.antrianRepo.GetByID(req.AntrianID)
	if err == nil {
		antrian.Status = "Selesai"
		s.antrianRepo.Update(antrian.ID, antrian)
	}

	return model.ToPemeriksaanResponse(createdPemeriksaan), nil
}

func (s *PemeriksaanService) GetPemeriksaanByID(ctx context.Context, id int) (model.PemeriksaanResponse, error) {
	pemeriksaan, err := s.repo.GetById(id)
	if err != nil {
		return model.PemeriksaanResponse{}, err
	}
	return model.ToPemeriksaanResponse(pemeriksaan), nil
}

func (s *PemeriksaanService) GetRiwayatPemeriksaanPasien(ctx context.Context, pasienID int) ([]model.PemeriksaanResponse, error) {
	allPemeriksaan, err := s.repo.GetAllByPasienID(pasienID)
	if err != nil {
		return nil, err
	}

	return model.ToPemeriksaanResponseList(allPemeriksaan), nil
}

func (s *PemeriksaanService) UpdatePemeriksaan(ctx context.Context, id int, req model.UpdatePemeriksaanRequest) (model.PemeriksaanResponse, error) {
	pemeriksaanUpdate := req.ToModel()

	updatedPemeriksaan, err := s.repo.Update(id, pemeriksaanUpdate)
	if err != nil {
		return model.PemeriksaanResponse{}, err
	}

	return model.ToPemeriksaanResponse(updatedPemeriksaan), nil
}

func (s *PemeriksaanService) DeletePemeriksaan(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
