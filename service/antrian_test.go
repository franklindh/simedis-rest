package service

import (
	"context"
	"errors"
	"testing"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAntrianService_CreateAntrian(t *testing.T) {

}

func TestAntrianService_GetAllAntrian(t *testing.T) {
	mockAntrianRepo := new(MockAntrianRepository)
	mockJadwalRepo := new(MockJadwalRepository)
	service := NewAntrianService(mockAntrianRepo, mockJadwalRepo)

	params := repository.ParamsGetAllAntrian{Page: 1, PageSize: 5}

	t.Run("Success: Get all antrian", func(t *testing.T) {
		mockAntrians := []model.Antrian{
			{ID: 1, Pasien: model.Pasien{NamaPasien: "Pasien A"}, Jadwal: model.Jadwal{Poli: model.Poli{Nama: "Poli A"}}},
			{ID: 2, Pasien: model.Pasien{NamaPasien: "Pasien B"}, Jadwal: model.Jadwal{Poli: model.Poli{Nama: "Poli B"}}},
		}
		mockMeta := pagination.Metadata{TotalRecords: 2}
		mockAntrianRepo.On("GetAll", params).Return(mockAntrians, mockMeta, nil).Once()

		results, meta, err := service.GetAllAntrian(context.Background(), params)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "Pasien B", results[1].Pasien.Nama)
		assert.Equal(t, 2, meta.TotalRecords)
		mockAntrianRepo.AssertExpectations(t)
	})
}

func TestAntrianService_GetAntrianByID(t *testing.T) {
	mockAntrianRepo := new(MockAntrianRepository)
	mockJadwalRepo := new(MockJadwalRepository)
	service := NewAntrianService(mockAntrianRepo, mockJadwalRepo)

	t.Run("Success: Antrian found", func(t *testing.T) {
		mockAntrian := model.Antrian{ID: 1, NomorAntrian: "G1", Pasien: model.Pasien{NamaPasien: "Pasien A"}}
		mockAntrianRepo.On("GetByID", 1).Return(mockAntrian, nil).Once()

		result, err := service.GetAntrianByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "G1", result.NomorAntrian)
		mockAntrianRepo.AssertExpectations(t)
	})

	t.Run("Fail: Antrian not found", func(t *testing.T) {
		mockAntrianRepo.On("GetByID", 99).Return(model.Antrian{}, repository.ErrNotFound).Once()

		_, err := service.GetAntrianByID(context.Background(), 99)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockAntrianRepo.AssertExpectations(t)
	})
}

func TestAntrianService_UpdateAntrian(t *testing.T) {
	mockAntrianRepo := new(MockAntrianRepository)
	mockJadwalRepo := new(MockJadwalRepository)
	service := NewAntrianService(mockAntrianRepo, mockJadwalRepo)

	req := model.UpdateAntrianRequest{Status: "Selesai", Prioritas: "Gawat"}

	t.Run("Success: Update antrian", func(t *testing.T) {
		updatedModel := req.ToModel()
		updatedModel.ID = 1
		mockAntrianRepo.On("Update", 1, mock.AnythingOfType("model.Antrian")).Return(updatedModel, nil).Once()

		result, err := service.UpdateAntrian(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "Selesai", result.Status)
		assert.Equal(t, "Gawat", result.Prioritas)
		mockAntrianRepo.AssertExpectations(t)
	})

	t.Run("Fail: Antrian to update not found", func(t *testing.T) {
		mockAntrianRepo.On("Update", 99, mock.AnythingOfType("model.Antrian")).Return(model.Antrian{}, repository.ErrNotFound).Once()

		_, err := service.UpdateAntrian(context.Background(), 99, req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockAntrianRepo.AssertExpectations(t)
	})
}

func TestAntrianService_DeleteAntrian(t *testing.T) {
	mockAntrianRepo := new(MockAntrianRepository)
	mockJadwalRepo := new(MockJadwalRepository)
	service := NewAntrianService(mockAntrianRepo, mockJadwalRepo)

	t.Run("Success: Delete antrian", func(t *testing.T) {
		mockAntrianRepo.On("Delete", 1).Return(nil).Once()
		err := service.DeleteAntrian(context.Background(), 1)
		assert.NoError(t, err)
		mockAntrianRepo.AssertExpectations(t)
	})

	t.Run("Fail: Antrian to delete not found", func(t *testing.T) {
		mockAntrianRepo.On("Delete", 99).Return(repository.ErrNotFound).Once()
		err := service.DeleteAntrian(context.Background(), 99)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockAntrianRepo.AssertExpectations(t)
	})
}
