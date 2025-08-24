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

type MockIcdRepository struct {
	mock.Mock
}

var _ IcdRepository = (*MockIcdRepository)(nil)

func (m *MockIcdRepository) Create(icd model.Icd) (model.Icd, error) {
	args := m.Called(icd)
	return args.Get(0).(model.Icd), args.Error(1)
}
func (m *MockIcdRepository) GetAll(params repository.ParamsGetAllIcd) ([]model.Icd, pagination.Metadata, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(pagination.Metadata), args.Error(2)
	}
	return args.Get(0).([]model.Icd), args.Get(1).(pagination.Metadata), args.Error(2)
}
func (m *MockIcdRepository) GetByID(id int) (model.Icd, error) {
	args := m.Called(id)
	return args.Get(0).(model.Icd), args.Error(1)
}
func (m *MockIcdRepository) Update(id int, icd model.Icd) (model.Icd, error) {
	args := m.Called(id, icd)
	return args.Get(0).(model.Icd), args.Error(1)
}
func (m *MockIcdRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestIcdService_CreateIcd(t *testing.T) {
	mockRepo := new(MockIcdRepository)
	service := NewIcdService(mockRepo)
	req := model.CreateIcdRequest{KodeIcd: "A01", NamaPenyakit: "Demam Tifoid", Status: "aktif"}

	t.Run("Success: Create ICD", func(t *testing.T) {
		createdModel := req.ToModel()
		createdModel.ID = 1
		mockRepo.On("Create", mock.AnythingOfType("model.Icd")).Return(createdModel, nil).Once()

		result, err := service.CreateIcd(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, "A01", result.KodeIcd)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: ICD code conflict", func(t *testing.T) {
		pgErr := &pgconn.PgError{Code: "23505"}
		mockRepo.On("Create", mock.AnythingOfType("model.Icd")).Return(model.Icd{}, pgErr).Once()

		_, err := service.CreateIcd(context.Background(), req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrIcdConflict))
		mockRepo.AssertExpectations(t)
	})
}

func TestIcdService_GetAllIcd(t *testing.T) {
	mockRepo := new(MockIcdRepository)
	service := NewIcdService(mockRepo)
	params := repository.ParamsGetAllIcd{Page: 1, PageSize: 5}

	t.Run("Success: Get all icd", func(t *testing.T) {
		mockIcds := []model.Icd{
			{ID: 1, KodeIcd: "A01", NamaPenyakit: "Demam Tifoid"},
			{ID: 2, KodeIcd: "J06", NamaPenyakit: "Infeksi Saluran Pernapasan Akut"},
		}
		mockMeta := pagination.Metadata{TotalRecords: 2}
		mockRepo.On("GetAll", params).Return(mockIcds, mockMeta, nil).Once()

		results, meta, err := service.GetAllIcd(context.Background(), params)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "J06", results[1].KodeIcd)
		assert.Equal(t, 2, meta.TotalRecords)
		mockRepo.AssertExpectations(t)
	})
}

func TestIcdService_GetIcdByID(t *testing.T) {
	mockRepo := new(MockIcdRepository)
	service := NewIcdService(mockRepo)

	t.Run("Success: ICD found", func(t *testing.T) {
		mockModel := model.Icd{ID: 1, KodeIcd: "A01", NamaPenyakit: "Demam Tifoid"}
		mockRepo.On("GetByID", 1).Return(mockModel, nil).Once()

		result, err := service.GetIcdByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "A01", result.KodeIcd)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: ICD not found", func(t *testing.T) {
		mockRepo.On("GetByID", 99).Return(model.Icd{}, repository.ErrNotFound).Once()

		_, err := service.GetIcdByID(context.Background(), 99)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}

func TestIcdService_UpdateIcd(t *testing.T) {
	mockRepo := new(MockIcdRepository)
	service := NewIcdService(mockRepo)
	req := model.UpdateIcdRequest{KodeIcd: "A01.1", NamaPenyakit: "Paratifoid A"}

	t.Run("Success: Update ICD", func(t *testing.T) {
		updatedModel := req.ToModel()
		updatedModel.ID = 1
		mockRepo.On("Update", 1, mock.AnythingOfType("model.Icd")).Return(updatedModel, nil).Once()

		result, err := service.UpdateIcd(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "A01.1", result.KodeIcd)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: ICD to update not found", func(t *testing.T) {
		mockRepo.On("Update", 99, mock.AnythingOfType("model.Icd")).Return(model.Icd{}, repository.ErrNotFound).Once()

		_, err := service.UpdateIcd(context.Background(), 99, req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}

func TestIcdService_DeleteIcd(t *testing.T) {
	mockRepo := new(MockIcdRepository)
	service := NewIcdService(mockRepo)

	t.Run("Success: Delete ICD", func(t *testing.T) {
		mockRepo.On("Delete", 1).Return(nil).Once()
		err := service.DeleteIcd(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: ICD to delete not found", func(t *testing.T) {
		mockRepo.On("Delete", 99).Return(repository.ErrNotFound).Once()
		err := service.DeleteIcd(context.Background(), 99)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}
