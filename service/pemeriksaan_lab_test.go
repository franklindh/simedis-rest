package service

import (
	"context"
	"errors"
	"testing"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPemeriksaanLabRepository struct {
	mock.Mock
}

var _ PemeriksaanLabRepository = (*MockPemeriksaanLabRepository)(nil)

func (m *MockPemeriksaanLabRepository) Create(hasilLab model.PemeriksaanLab) (model.PemeriksaanLab, error) {
	args := m.Called(hasilLab)

	if retFn, ok := args.Get(0).(func(model.PemeriksaanLab) model.PemeriksaanLab); ok {
		return retFn(hasilLab), args.Error(1)
	}

	return args.Get(0).(model.PemeriksaanLab), args.Error(1)
}
func (m *MockPemeriksaanLabRepository) GetAllByPemeriksaanID(pemeriksaanID int) ([]model.PemeriksaanLab, error) {
	args := m.Called(pemeriksaanID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PemeriksaanLab), args.Error(1)
}
func (m *MockPemeriksaanLabRepository) GetById(id int) (model.PemeriksaanLab, error) {
	args := m.Called(id)
	return args.Get(0).(model.PemeriksaanLab), args.Error(1)
}
func (m *MockPemeriksaanLabRepository) Update(id int, hasilLab model.PemeriksaanLab) (model.PemeriksaanLab, error) {
	args := m.Called(id, hasilLab)
	return args.Get(0).(model.PemeriksaanLab), args.Error(1)
}
func (m *MockPemeriksaanLabRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestPemeriksaanLabService_GetAllByPemeriksaanID(t *testing.T) {
	mockRepo := new(MockPemeriksaanLabRepository)
	service := NewPemeriksaanLabService(mockRepo)
	pemeriksaanID := 100

	t.Run("Success: Get all lab results for a pemeriksaan", func(t *testing.T) {

		mockResults := []model.PemeriksaanLab{
			{ID: 1, Hasil: "150.000", JenisPemeriksaanLab: model.JenisPemeriksaanLab{ID: 1, NamaPemeriksaan: "Trombosit"}},
			{ID: 2, Hasil: "Negatif", JenisPemeriksaanLab: model.JenisPemeriksaanLab{ID: 2, NamaPemeriksaan: "Urine"}},
		}
		mockRepo.On("GetAllByPemeriksaanID", pemeriksaanID).Return(mockResults, nil).Once()

		results, err := service.GetAllByPemeriksaanID(context.Background(), pemeriksaanID)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "Trombosit", results[0].JenisPemeriksaanLab.NamaPemeriksaan)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success: Return empty slice when no results found", func(t *testing.T) {

		mockRepo.On("GetAllByPemeriksaanID", pemeriksaanID).Return([]model.PemeriksaanLab{}, nil).Once()

		results, err := service.GetAllByPemeriksaanID(context.Background(), pemeriksaanID)

		assert.NoError(t, err)
		assert.Len(t, results, 0)
		mockRepo.AssertExpectations(t)
	})
}

func TestPemeriksaanLabService_CreateBatch(t *testing.T) {
	mockRepo := new(MockPemeriksaanLabRepository)
	service := NewPemeriksaanLabService(mockRepo)

	reqs := []model.CreateHasilLabRequest{
		{JenisPemeriksaanID: 1, Hasil: "150.000"},
		{JenisPemeriksaanID: 2, Hasil: "Negatif"},
	}
	pemeriksaanID := 100

	t.Run("Success: Create batch", func(t *testing.T) {

		mockRepo.On("Create", mock.AnythingOfType("model.PemeriksaanLab")).
			Return(func(lab model.PemeriksaanLab) model.PemeriksaanLab {

				lab.ID = lab.JenisPemeriksaanID
				return lab
			}, nil).Times(2)

		results, err := service.CreateBatch(context.Background(), pemeriksaanID, reqs)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "150.000", results[0].Hasil)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: One of the batch items fails", func(t *testing.T) {

		mockRepo.On("Create", mock.AnythingOfType("model.PemeriksaanLab")).Return(model.PemeriksaanLab{ID: 1}, nil).Once()
		mockRepo.On("Create", mock.AnythingOfType("model.PemeriksaanLab")).Return(model.PemeriksaanLab{}, errors.New("db error")).Once()

		_, err := service.CreateBatch(context.Background(), pemeriksaanID, reqs)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		mockRepo.AssertExpectations(t)
	})
}

func TestPemeriksaanLabService_Update(t *testing.T) {
	mockRepo := new(MockPemeriksaanLabRepository)
	service := NewPemeriksaanLabService(mockRepo)

	req := model.UpdateHasilLabRequest{Hasil: "Positif"}

	t.Run("Success: Update lab result", func(t *testing.T) {
		mockModel := req.ToModel()
		mockModel.ID = 1
		mockRepo.On("Update", 1, mock.AnythingOfType("model.PemeriksaanLab")).Return(mockModel, nil).Once()

		result, err := service.Update(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.Equal(t, "Positif", result.Hasil)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Lab result not found", func(t *testing.T) {
		mockRepo.On("Update", 99, mock.AnythingOfType("model.PemeriksaanLab")).Return(model.PemeriksaanLab{}, repository.ErrNotFound).Once()
		_, err := service.Update(context.Background(), 99, req)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}

func TestPemeriksaanLabService_Delete(t *testing.T) {
	mockRepo := new(MockPemeriksaanLabRepository)
	service := NewPemeriksaanLabService(mockRepo)

	t.Run("Success: Delete lab result", func(t *testing.T) {

		mockRepo.On("Delete", 1).Return(nil).Once()

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Lab result to delete not found", func(t *testing.T) {

		mockRepo.On("Delete", 99).Return(repository.ErrNotFound).Once()

		err := service.Delete(context.Background(), 99)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}
