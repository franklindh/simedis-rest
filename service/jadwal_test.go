package service

import (
	"context"
	"errors"
	"testing"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestJadwalService_CreateJadwal(t *testing.T) {
	t.Run("Success: Create new schedule", func(t *testing.T) {
		mockRepo := new(MockJadwalRepository)
		service := NewJadwalService(mockRepo)
		req := model.JadwalRequest{PetugasID: 1, PoliID: 1, Tanggal: "2025-08-23", WaktuMulai: "09:00", WaktuSelesai: "11:00"}

		createdModel := req.ToModel()
		createdModel.ID = 10
		fullModel := createdModel
		fullModel.Petugas = model.Petugas{ID: 1, Nama: "Dr. Budi"}
		fullModel.Poli = model.Poli{ID: 1, Nama: "Poli Umum"}

		mockRepo.On("Create", mock.AnythingOfType("model.Jadwal")).Return(createdModel, nil).Once()
		mockRepo.On("GetById", 10).Return(fullModel, nil).Once()

		result, err := service.CreateJadwal(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, 10, result.ID)
		assert.Equal(t, "Dr. Budi", result.Petugas.Nama)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Invalid end time", func(t *testing.T) {
		mockRepo := new(MockJadwalRepository)
		service := NewJadwalService(mockRepo)
		invalidReq := model.JadwalRequest{WaktuMulai: "11:00", WaktuSelesai: "09:00"}

		_, err := service.CreateJadwal(context.Background(), invalidReq)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrWaktuSelesaiInvalid))
		mockRepo.AssertNotCalled(t, "Create", mock.Anything)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Schedule conflict", func(t *testing.T) {
		mockRepo := new(MockJadwalRepository)
		service := NewJadwalService(mockRepo)
		req := model.JadwalRequest{PetugasID: 1, PoliID: 1, Tanggal: "2025-08-23", WaktuMulai: "09:00", WaktuSelesai: "11:00"}

		pgErr := &pgconn.PgError{Code: "23505"}
		mockRepo.On("Create", mock.AnythingOfType("model.Jadwal")).Return(model.Jadwal{}, pgErr).Once()

		_, err := service.CreateJadwal(context.Background(), req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrJadwalConflict))
		mockRepo.AssertExpectations(t)
	})
}

func TestJadwalService_GetAllJadwal(t *testing.T) {
	mockRepo := new(MockJadwalRepository)
	service := NewJadwalService(mockRepo)
	params := repository.ParamsGetAllJadwal{Page: 1, PageSize: 5}

	t.Run("Success: Get all jadwal", func(t *testing.T) {
		mockJadwals := []model.Jadwal{
			{ID: 1, Petugas: model.Petugas{Nama: "Dr. Budi"}, Poli: model.Poli{Nama: "Poli Umum"}},
			{ID: 2, Petugas: model.Petugas{Nama: "Dr. Ani"}, Poli: model.Poli{Nama: "Poli Gigi"}},
		}
		mockMeta := pagination.Metadata{TotalRecords: 2}
		mockRepo.On("GetAll", params).Return(mockJadwals, mockMeta, nil).Once()

		results, meta, err := service.GetAllJadwal(context.Background(), params)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "Dr. Ani", results[1].Petugas.Nama)
		assert.Equal(t, 2, meta.TotalRecords)
		mockRepo.AssertExpectations(t)
	})
}

func TestJadwalService_GetJadwalByID(t *testing.T) {
	mockRepo := new(MockJadwalRepository)
	service := NewJadwalService(mockRepo)

	t.Run("Success: Jadwal found", func(t *testing.T) {
		fullModel := model.Jadwal{
			ID: 1, Petugas: model.Petugas{ID: 1, Nama: "Dr. Budi"}, Poli: model.Poli{ID: 1, Nama: "Poli Umum"},
		}
		mockRepo.On("GetById", 1).Return(fullModel, nil).Once()

		result, err := service.GetJadwalByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "Dr. Budi", result.Petugas.Nama)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Jadwal not found", func(t *testing.T) {
		mockRepo.On("GetById", 99).Return(model.Jadwal{}, repository.ErrNotFound).Once()
		_, err := service.GetJadwalByID(context.Background(), 99)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}

func TestJadwalService_UpdateJadwal(t *testing.T) {
	mockRepo := new(MockJadwalRepository)
	service := NewJadwalService(mockRepo)
	req := model.JadwalRequest{
		PetugasID: 1, PoliID: 1, Tanggal: "2025-08-24", WaktuMulai: "13:00", WaktuSelesai: "15:00",
	}

	t.Run("Success: Update jadwal", func(t *testing.T) {
		updatedModel := req.ToModel()
		updatedModel.ID = 1

		fullModel := updatedModel
		fullModel.Petugas = model.Petugas{ID: 1, Nama: "Dr. Budi Updated"}
		fullModel.Poli = model.Poli{ID: 1, Nama: "Poli Umum Updated"}

		mockRepo.On("Update", 1, mock.AnythingOfType("model.Jadwal")).Return(updatedModel, nil).Once()

		mockRepo.On("GetById", 1).Return(fullModel, nil).Once()

		result, err := service.UpdateJadwal(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "Dr. Budi Updated", result.Petugas.Nama)
		mockRepo.AssertExpectations(t)
	})
}

func TestJadwalService_DeleteJadwal(t *testing.T) {
	mockRepo := new(MockJadwalRepository)
	service := NewJadwalService(mockRepo)

	t.Run("Success: Delete jadwal", func(t *testing.T) {
		mockRepo.On("Delete", 1).Return(nil).Once()
		err := service.DeleteJadwal(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Jadwal to delete not found", func(t *testing.T) {
		mockRepo.On("Delete", 99).Return(repository.ErrNotFound).Once()
		err := service.DeleteJadwal(context.Background(), 99)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}
