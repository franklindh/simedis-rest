package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPemeriksaanService_CreatePemeriksaan(t *testing.T) {
	mockPemeriksaanRepo := new(MockPemeriksaanRepository)
	mockAntrianRepo := new(MockAntrianRepository)
	pemeriksaanService := NewPemeriksaanService(mockPemeriksaanRepo, mockAntrianRepo)

	inputDTO := model.CreatePemeriksaanRequest{
		AntrianID:          1,
		TanggalPemeriksaan: "2025-08-21",
		Keluhan:            "Pusing",
	}

	t.Run("Success: Create pemeriksaan and update antrian status", func(t *testing.T) {
		createdModel := inputDTO.ToModel()
		createdModel.ID = 10
		mockPemeriksaanRepo.On("CheckExistingPemeriksaan", 1).Return(repository.ErrNotFound).Once()
		mockPemeriksaanRepo.On("Create", mock.AnythingOfType("model.Pemeriksaan")).Return(createdModel, nil).Once()
		mockAntrian := model.Antrian{ID: 1, Status: "Menunggu Diagnosis"}
		mockAntrianRepo.On("GetByID", 1).Return(mockAntrian, nil).Once()
		mockAntrianRepo.On("Update", 1, mock.AnythingOfType("model.Antrian")).Return(model.Antrian{}, nil).Once()

		result, err := pemeriksaanService.CreatePemeriksaan(context.Background(), inputDTO)

		assert.NoError(t, err)
		assert.Equal(t, 10, result.ID)
		mockPemeriksaanRepo.AssertExpectations(t)
		mockAntrianRepo.AssertExpectations(t)
		mockAntrianRepo.AssertCalled(t, "Update", 1, mock.Anything)
	})

	t.Run("Fail: Pemeriksaan already exists", func(t *testing.T) {
		mockPemeriksaanRepo.On("CheckExistingPemeriksaan", 1).Return(nil).Once()
		_, err := pemeriksaanService.CreatePemeriksaan(context.Background(), inputDTO)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrPemeriksaanExists))
		mockPemeriksaanRepo.AssertExpectations(t)
	})
}

func TestPemeriksaanService_GetPemeriksaanByID(t *testing.T) {
	mockPemeriksaanRepo := new(MockPemeriksaanRepository)
	mockAntrianRepo := new(MockAntrianRepository)
	pemeriksaanService := NewPemeriksaanService(mockPemeriksaanRepo, mockAntrianRepo)

	t.Run("Success: Pemeriksaan found", func(t *testing.T) {
		mockModel := model.Pemeriksaan{
			ID:                 1,
			TanggalPemeriksaan: time.Now(),
			Keluhan:            sql.NullString{String: "Sakit Perut", Valid: true},
			Antrian: model.Antrian{
				Pasien: model.Pasien{ID: 5, NamaPasien: "Budi"},
				Jadwal: model.Jadwal{
					Petugas: model.Petugas{ID: 2, Nama: "Dr. Ani"},
					Poli:    model.Poli{ID: 3, Nama: "Poli Umum"},
				},
			},
		}
		mockPemeriksaanRepo.On("GetById", 1).Return(mockModel, nil).Once()

		result, err := pemeriksaanService.GetPemeriksaanByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "Budi", result.Pasien.Nama)
		assert.Equal(t, "Dr. Ani", result.Dokter.Nama)
		mockPemeriksaanRepo.AssertExpectations(t)
	})

	t.Run("Fail: Pemeriksaan not found", func(t *testing.T) {
		mockPemeriksaanRepo.On("GetById", 99).Return(model.Pemeriksaan{}, repository.ErrNotFound).Once()
		_, err := pemeriksaanService.GetPemeriksaanByID(context.Background(), 99)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockPemeriksaanRepo.AssertExpectations(t)
	})
}

func TestPemeriksaanService_GetRiwayatPemeriksaanPasien(t *testing.T) {
	mockPemeriksaanRepo := new(MockPemeriksaanRepository)
	mockAntrianRepo := new(MockAntrianRepository)
	pemeriksaanService := NewPemeriksaanService(mockPemeriksaanRepo, mockAntrianRepo)

	t.Run("Success: Get patient history", func(t *testing.T) {
		mockHistory := []model.Pemeriksaan{
			{
				ID: 1, TanggalPemeriksaan: time.Now(), Keluhan: sql.NullString{String: "Batuk", Valid: true},
				Antrian: model.Antrian{Pasien: model.Pasien{ID: 5, NamaPasien: "Andi"}, Jadwal: model.Jadwal{Petugas: model.Petugas{ID: 2, Nama: "Dr. Ani"}, Poli: model.Poli{ID: 3, Nama: "Poli Umum"}}},
			},
			{
				ID: 2, TanggalPemeriksaan: time.Now().AddDate(0, -1, 0), Keluhan: sql.NullString{String: "Pilek", Valid: true},
				Antrian: model.Antrian{Pasien: model.Pasien{ID: 5, NamaPasien: "Andi"}, Jadwal: model.Jadwal{Petugas: model.Petugas{ID: 2, Nama: "Dr. Ani"}, Poli: model.Poli{ID: 3, Nama: "Poli Umum"}}},
			},
		}
		mockPemeriksaanRepo.On("GetAllByPasienID", 5).Return(mockHistory, nil).Once()

		result, err := pemeriksaanService.GetRiwayatPemeriksaanPasien(context.Background(), 5)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "Batuk", result[0].Keluhan)
		assert.Equal(t, "Andi", result[0].Pasien.Nama)
		mockPemeriksaanRepo.AssertExpectations(t)
	})
}

func TestPemeriksaanService_UpdatePemeriksaan(t *testing.T) {
	mockPemeriksaanRepo := new(MockPemeriksaanRepository)
	mockAntrianRepo := new(MockAntrianRepository)
	pemeriksaanService := NewPemeriksaanService(mockPemeriksaanRepo, mockAntrianRepo)

	inputDTO := model.UpdatePemeriksaanRequest{
		TanggalPemeriksaan: "2025-08-22",
		Keluhan:            "Sudah mendingan",
	}

	t.Run("Success: Update pemeriksaan", func(t *testing.T) {
		updatedModel := inputDTO.ToModel()
		updatedModel.ID = 1
		mockPemeriksaanRepo.On("Update", 1, mock.AnythingOfType("model.Pemeriksaan")).Return(updatedModel, nil).Once()

		result, err := pemeriksaanService.UpdatePemeriksaan(context.Background(), 1, inputDTO)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "Sudah mendingan", result.Keluhan)
		mockPemeriksaanRepo.AssertExpectations(t)
	})

	t.Run("Fail: Pemeriksaan to update not found", func(t *testing.T) {
		mockPemeriksaanRepo.On("Update", 99, mock.AnythingOfType("model.Pemeriksaan")).Return(model.Pemeriksaan{}, repository.ErrNotFound).Once()

		_, err := pemeriksaanService.UpdatePemeriksaan(context.Background(), 99, inputDTO)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockPemeriksaanRepo.AssertExpectations(t)
	})
}

func TestPemeriksaanService_DeletePemeriksaan(t *testing.T) {
	mockPemeriksaanRepo := new(MockPemeriksaanRepository)
	mockAntrianRepo := new(MockAntrianRepository)
	pemeriksaanService := NewPemeriksaanService(mockPemeriksaanRepo, mockAntrianRepo)

	t.Run("Success: Delete pemeriksaan", func(t *testing.T) {
		mockPemeriksaanRepo.On("Delete", 1).Return(nil).Once()
		err := pemeriksaanService.DeletePemeriksaan(context.Background(), 1)
		assert.NoError(t, err)
		mockPemeriksaanRepo.AssertExpectations(t)
	})

	t.Run("Fail: Pemeriksaan to delete not found", func(t *testing.T) {
		mockPemeriksaanRepo.On("Delete", 99).Return(repository.ErrNotFound).Once()
		err := pemeriksaanService.DeletePemeriksaan(context.Background(), 99)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockPemeriksaanRepo.AssertExpectations(t)
	})
}
