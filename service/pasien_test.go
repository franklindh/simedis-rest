package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPasienRepository struct {
	mock.Mock
}

var _ PasienRepository = (*MockPasienRepository)(nil)

func (m *MockPasienRepository) Create(pasien model.Pasien) (model.Pasien, error) {
	args := m.Called(pasien)

	if retFn, ok := args.Get(0).(func(model.Pasien) model.Pasien); ok {
		return retFn(pasien), args.Error(1)
	}
	return args.Get(0).(model.Pasien), args.Error(1)
}

func (m *MockPasienRepository) GetAll(params repository.ParamsGetAllPasien) ([]model.Pasien, pagination.Metadata, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(pagination.Metadata), args.Error(2)
	}
	return args.Get(0).([]model.Pasien), args.Get(1).(pagination.Metadata), args.Error(2)
}

func (m *MockPasienRepository) GetById(id int) (model.Pasien, error) {
	args := m.Called(id)
	return args.Get(0).(model.Pasien), args.Error(1)
}

func (m *MockPasienRepository) Update(id int, pasien model.Pasien) (model.Pasien, error) {
	args := m.Called(id, pasien)
	return args.Get(0).(model.Pasien), args.Error(1)
}

func (m *MockPasienRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPasienRepository) GetLastID() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func TestPasienService_CreatePasien(t *testing.T) {
	mockRepo := new(MockPasienRepository)
	service := NewPasienService(mockRepo)

	req := model.CreatePasienRequest{
		NIK:                "1234567890123456",
		NamaPasien:         "Budi",
		AlamatPasien:       "Jl. Sehat",
		TempatLahirPasien:  "Jakarta",
		TanggalLahirPasien: "2000-01-01",
		JKPasien:           "L",
		StatusPernikahan:   "Belum Menikah",
	}

	t.Run("Success: Create new patient", func(t *testing.T) {

		mockRepo.On("GetLastID").Return(9, nil).Once()
		mockRepo.On("Create", mock.AnythingOfType("model.Pasien")).
			Return(func(p model.Pasien) model.Pasien {
				p.ID = 10
				return p
			}, nil).Once()

		result, err := service.CreatePasien(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, 10, result.ID)
		assert.Equal(t, "Budi", result.NamaPasien)

		expectedSuffix := "-0010"
		assert.True(t, strings.HasPrefix(result.NoRekamMedis, "RM-"))
		assert.True(t, strings.HasSuffix(result.NoRekamMedis, expectedSuffix))

		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: NIK conflict", func(t *testing.T) {
		mockRepo.On("GetLastID").Return(9, nil).Once()
		pgErr := &pgconn.PgError{Code: "23505"}
		mockRepo.On("Create", mock.AnythingOfType("model.Pasien")).Return(model.Pasien{}, pgErr).Once()

		_, err := service.CreatePasien(context.Background(), req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrPasienConflict))
		mockRepo.AssertExpectations(t)
	})
}

func TestPasienService_GetAllPasien(t *testing.T) {
	mockRepo := new(MockPasienRepository)
	service := NewPasienService(mockRepo)
	params := repository.ParamsGetAllPasien{Page: 1, PageSize: 5}

	t.Run("Success: Get all pasien", func(t *testing.T) {
		mockPasiens := []model.Pasien{
			{ID: 1, NamaPasien: "Andi"},
			{ID: 2, NamaPasien: "Budi"},
		}
		mockMeta := pagination.Metadata{CurrentPage: 1, PageSize: 5, TotalRecords: 2}
		mockRepo.On("GetAll", params).Return(mockPasiens, mockMeta, nil).Once()

		results, meta, err := service.GetAllPasien(context.Background(), params)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "Andi", results[0].NamaPasien)
		assert.Equal(t, 2, meta.TotalRecords)
		mockRepo.AssertExpectations(t)
	})
}

func TestPasienService_GetPasienByID(t *testing.T) {
	mockRepo := new(MockPasienRepository)
	service := NewPasienService(mockRepo)

	t.Run("Success: Pasien found", func(t *testing.T) {
		mockPasien := model.Pasien{ID: 1, NamaPasien: "Cici"}
		mockRepo.On("GetById", 1).Return(mockPasien, nil).Once()

		result, err := service.GetPasienByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "Cici", result.NamaPasien)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Pasien not found", func(t *testing.T) {
		mockRepo.On("GetById", 99).Return(model.Pasien{}, repository.ErrNotFound).Once()

		_, err := service.GetPasienByID(context.Background(), 99)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}

func TestPasienService_UpdatePasien(t *testing.T) {
	mockRepo := new(MockPasienRepository)
	service := NewPasienService(mockRepo)

	req := model.UpdatePasienRequest{
		NIK:        "1234567890123456",
		NamaPasien: "Budi Updated",
	}

	t.Run("Success: Update pasien", func(t *testing.T) {
		updatedModel := req.ToModel()
		updatedModel.ID = 1
		mockRepo.On("Update", 1, mock.AnythingOfType("model.Pasien")).Return(updatedModel, nil).Once()

		result, err := service.UpdatePasien(context.Background(), 1, req)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "Budi Updated", result.NamaPasien)
		mockRepo.AssertExpectations(t)
	})
}

func TestPasienService_DeletePasien(t *testing.T) {
	mockRepo := new(MockPasienRepository)
	service := NewPasienService(mockRepo)

	t.Run("Success: Delete pasien", func(t *testing.T) {
		mockRepo.On("Delete", 1).Return(nil).Once()
		err := service.DeletePasien(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Pasien to delete not found", func(t *testing.T) {
		mockRepo.On("Delete", 99).Return(repository.ErrNotFound).Once()
		err := service.DeletePasien(context.Background(), 99)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}
