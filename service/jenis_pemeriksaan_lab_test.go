package service

import (
	"context"
	"errors"
	"testing"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockJenisPemeriksaanLabRepository struct {
	mock.Mock
}

var _ JenisPemeriksaanLabRepository = (*MockJenisPemeriksaanLabRepository)(nil)

func (m *MockJenisPemeriksaanLabRepository) Create(jenis model.JenisPemeriksaanLab) (model.JenisPemeriksaanLab, error) {
	args := m.Called(jenis)
	return args.Get(0).(model.JenisPemeriksaanLab), args.Error(1)
}
func (m *MockJenisPemeriksaanLabRepository) GetAll() ([]model.JenisPemeriksaanLab, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.JenisPemeriksaanLab), args.Error(1)
}
func (m *MockJenisPemeriksaanLabRepository) GetById(id int) (model.JenisPemeriksaanLab, error) {
	args := m.Called(id)
	return args.Get(0).(model.JenisPemeriksaanLab), args.Error(1)
}
func (m *MockJenisPemeriksaanLabRepository) Update(id int, jenis model.JenisPemeriksaanLab) (model.JenisPemeriksaanLab, error) {
	args := m.Called(id, jenis)
	return args.Get(0).(model.JenisPemeriksaanLab), args.Error(1)
}
func (m *MockJenisPemeriksaanLabRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockJenisPemeriksaanLabRepository) FindByName(name string) (model.JenisPemeriksaanLab, error) {
	args := m.Called(name)
	return args.Get(0).(model.JenisPemeriksaanLab), args.Error(1)
}

func TestJenisPemeriksaanLabService_Create(t *testing.T) {
	mockRepo := new(MockJenisPemeriksaanLabRepository)
	service := NewJenisPemeriksaanLabService(mockRepo)
	req := model.CreateJenisPemeriksaanLabRequest{NamaPemeriksaan: "Leukosit"}

	t.Run("Success: Create", func(t *testing.T) {

		createdModel := req.ToModel()
		createdModel.ID = 1
		mockRepo.On("Create", mock.AnythingOfType("model.JenisPemeriksaanLab")).Return(createdModel, nil).Once()

		result, err := service.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "Leukosit", result.NamaPemeriksaan)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Name conflict", func(t *testing.T) {

		pgErr := &pgconn.PgError{Code: "23505"}
		mockRepo.On("Create", mock.AnythingOfType("model.JenisPemeriksaanLab")).Return(model.JenisPemeriksaanLab{}, pgErr).Once()

		_, err := service.Create(context.Background(), req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrJenisPemeriksaanConflict))
		mockRepo.AssertExpectations(t)
	})
}

func TestJenisPemeriksaanLabService_GetAll(t *testing.T) {
	mockRepo := new(MockJenisPemeriksaanLabRepository)
	service := NewJenisPemeriksaanLabService(mockRepo)

	t.Run("Success: Get all", func(t *testing.T) {
		mockList := []model.JenisPemeriksaanLab{
			{ID: 1, NamaPemeriksaan: "Leukosit"},
			{ID: 2, NamaPemeriksaan: "Trombosit"},
		}
		mockRepo.On("GetAll").Return(mockList, nil).Once()

		results, err := service.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "Trombosit", results[1].NamaPemeriksaan)
		mockRepo.AssertExpectations(t)
	})
}

func TestJenisPemeriksaanLabService_GetByID(t *testing.T) {
	mockRepo := new(MockJenisPemeriksaanLabRepository)
	service := NewJenisPemeriksaanLabService(mockRepo)

	t.Run("Success: Found", func(t *testing.T) {
		mockModel := model.JenisPemeriksaanLab{ID: 1, NamaPemeriksaan: "Leukosit"}
		mockRepo.On("GetById", 1).Return(mockModel, nil).Once()

		result, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "Leukosit", result.NamaPemeriksaan)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Not Found", func(t *testing.T) {
		mockRepo.On("GetById", 99).Return(model.JenisPemeriksaanLab{}, repository.ErrNotFound).Once()

		_, err := service.GetByID(context.Background(), 99)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}

func TestJenisPemeriksaanLabService_Update(t *testing.T) {
	mockRepo := new(MockJenisPemeriksaanLabRepository)
	service := NewJenisPemeriksaanLabService(mockRepo)
	req := model.UpdateJenisPemeriksaanLabRequest{NamaPemeriksaan: "Trombosit Baru"}

	t.Run("Success: Update", func(t *testing.T) {
		mockRepo.On("FindByName", "Trombosit Baru").Return(model.JenisPemeriksaanLab{}, repository.ErrNotFound).Once()
		updatedModel := req.ToModel()
		updatedModel.ID = 1
		mockRepo.On("Update", 1, mock.AnythingOfType("model.JenisPemeriksaanLab")).Return(updatedModel, nil).Once()

		result, err := service.Update(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.Equal(t, "Trombosit Baru", result.NamaPemeriksaan)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Name conflict with another item", func(t *testing.T) {
		existing := model.JenisPemeriksaanLab{ID: 2, NamaPemeriksaan: "Trombosit Baru"}
		mockRepo.On("FindByName", "Trombosit Baru").Return(existing, nil).Once()

		_, err := service.Update(context.Background(), 1, req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrJenisPemeriksaanConflict))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Item to update not found", func(t *testing.T) {
		mockRepo.On("FindByName", "Trombosit Baru").Return(model.JenisPemeriksaanLab{}, repository.ErrNotFound).Once()
		mockRepo.On("Update", 99, mock.AnythingOfType("model.JenisPemeriksaanLab")).Return(model.JenisPemeriksaanLab{}, repository.ErrNotFound).Once()

		_, err := service.Update(context.Background(), 99, req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}

func TestJenisPemeriksaanLabService_Delete(t *testing.T) {
	mockRepo := new(MockJenisPemeriksaanLabRepository)
	service := NewJenisPemeriksaanLabService(mockRepo)

	t.Run("Success: Delete", func(t *testing.T) {
		mockRepo.On("Delete", 1).Return(nil).Once()
		err := service.Delete(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Item to delete not found", func(t *testing.T) {
		mockRepo.On("Delete", 99).Return(repository.ErrNotFound).Once()
		err := service.Delete(context.Background(), 99)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}
