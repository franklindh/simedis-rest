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

type MockPoliRepository struct {
	mock.Mock
}

var _ repository.PoliRepositoryInterface = (*MockPoliRepository)(nil)

func (m *MockPoliRepository) Create(poli model.Poli) (model.Poli, error) {
	args := m.Called(poli)
	return args.Get(0).(model.Poli), args.Error(1)
}

func (m *MockPoliRepository) GetAll() ([]model.Poli, error) {
	args := m.Called()
	return args.Get(0).([]model.Poli), args.Error(1)
}

func (m *MockPoliRepository) GetByID(id int) (model.Poli, error) {
	args := m.Called(id)
	return args.Get(0).(model.Poli), args.Error(1)
}

func (m *MockPoliRepository) Update(id int, poli model.Poli) (model.Poli, error) {
	args := m.Called(id, poli)
	return args.Get(0).(model.Poli), args.Error(1)
}

func (m *MockPoliRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPoliRepository) FindByName(name string) (model.Poli, error) {
	args := m.Called(name)
	return args.Get(0).(model.Poli), args.Error(1)
}

func TestPoliService_CreatePoli(t *testing.T) {
	testCases := []struct {
		name          string
		inputDTO      model.CreatePoliRequest
		setupMock     func(mockRepo *MockPoliRepository, poli model.Poli)
		expectedError error
	}{
		{
			name:     "Success: Poli created successfully",
			inputDTO: model.CreatePoliRequest{Name: "Poli Jantung", Status: "aktif"},
			setupMock: func(mockRepo *MockPoliRepository, poli model.Poli) {
				expectedResult := poli
				expectedResult.ID = 1
				mockRepo.On("Create", poli).Return(expectedResult, nil).Once()
			},
			expectedError: nil,
		},
		{
			name:     "Fail: Name conflict",
			inputDTO: model.CreatePoliRequest{Name: "Poli Jantung", Status: "aktif"},
			setupMock: func(mockRepo *MockPoliRepository, poli model.Poli) {
				pgErr := &pgconn.PgError{Code: "23505"}
				mockRepo.On("Create", poli).Return(model.Poli{}, pgErr).Once()
			},
			expectedError: ErrPoliConflict,
		},
		{
			name:     "Fail: Generic database error",
			inputDTO: model.CreatePoliRequest{Name: "Poli Jantung", Status: "aktif"},
			setupMock: func(mockRepo *MockPoliRepository, poli model.Poli) {
				mockRepo.On("Create", poli).Return(model.Poli{}, errors.New("unexpected database error")).Once()
			},
			expectedError: errors.New("failed to create poli"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockPoliRepository)
			poliModel := tc.inputDTO.ToModel()
			tc.setupMock(mockRepo, poliModel)

			poliService := NewPoliService(mockRepo)
			_, err := poliService.CreatePoli(context.Background(), tc.inputDTO)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestPoliService_GetAllPolis(t *testing.T) {
	testCases := []struct {
		name          string
		setupMock     func(mockRepo *MockPoliRepository)
		expectedCount int
		expectedError error
	}{
		{
			name: "Success: Returns multiple poli",
			setupMock: func(mockRepo *MockPoliRepository) {
				expectedPolis := []model.Poli{
					{ID: 1, Nama: "Poli Gigi"},
					{ID: 2, Nama: "Poli Anak"},
				}
				mockRepo.On("GetAll").Return(expectedPolis, nil).Once()
			},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "Success: Returns empty slice when no data",
			setupMock: func(mockRepo *MockPoliRepository) {
				mockRepo.On("GetAll").Return([]model.Poli{}, nil).Once()
			},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name: "Fail: Database error",
			setupMock: func(mockRepo *MockPoliRepository) {
				mockRepo.On("GetAll").Return([]model.Poli{}, errors.New("unexpected database error")).Once()
			},
			expectedCount: 0,
			expectedError: errors.New("unexpected database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockPoliRepository)
			tc.setupMock(mockRepo)
			poliService := NewPoliService(mockRepo)

			result, err := poliService.GetAllPolis(context.Background())

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tc.expectedCount)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPoliService_GetPoliByID(t *testing.T) {
	testCases := []struct {
		name          string
		inputID       int
		setupMock     func(mockRepo *MockPoliRepository)
		expectedError error
	}{
		{
			name:    "Success: Poli found",
			inputID: 1,
			setupMock: func(mockRepo *MockPoliRepository) {
				expectedPoli := model.Poli{ID: 1, Nama: "Poli Umum", Status: "aktif"}
				mockRepo.On("GetByID", 1).Return(expectedPoli, nil).Once()
			},
			expectedError: nil,
		},
		{
			name:    "Fail: Poli not found",
			inputID: 99,
			setupMock: func(mockRepo *MockPoliRepository) {
				mockRepo.On("GetByID", 99).Return(model.Poli{}, repository.ErrNotFound).Once()
			},
			expectedError: repository.ErrNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockPoliRepository)
			tc.setupMock(mockRepo)
			poliService := NewPoliService(mockRepo)

			_, err := poliService.GetPoliByID(context.Background(), tc.inputID)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError))
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPoliService_UpdatePoli(t *testing.T) {

	testCases := []struct {
		name          string
		inputID       int
		inputDTO      model.UpdatePoliRequest
		setupMock     func(mockRepo *MockPoliRepository)
		expectedError error
	}{
		{
			name:     "Success: Poli updated successfully",
			inputID:  1,
			inputDTO: model.UpdatePoliRequest{Name: "Poli Lansia", Status: "aktif"},
			setupMock: func(mockRepo *MockPoliRepository) {
				updatedModel := model.Poli{ID: 1, Nama: "Poli Lansia", Status: "aktif"}
				mockRepo.On("FindByName", "Poli Lansia").Return(model.Poli{}, repository.ErrNotFound).Once()
				mockRepo.On("Update", 1, mock.AnythingOfType("model.Poli")).Return(updatedModel, nil).Once()
			},
			expectedError: nil,
		},
		{
			name:     "Fail: Poli not found",
			inputID:  99,
			inputDTO: model.UpdatePoliRequest{Name: "Poli Hantu", Status: "aktif"},
			setupMock: func(mockRepo *MockPoliRepository) {
				mockRepo.On("FindByName", "Poli Hantu").Return(model.Poli{}, repository.ErrNotFound).Once()
				mockRepo.On("Update", 99, mock.AnythingOfType("model.Poli")).Return(model.Poli{}, repository.ErrNotFound).Once()
			},
			expectedError: repository.ErrNotFound,
		},
		{
			name:     "Fail: Name conflict with another poli",
			inputID:  1,
			inputDTO: model.UpdatePoliRequest{Name: "Poli Gigi", Status: "aktif"},
			setupMock: func(mockRepo *MockPoliRepository) {
				existingPoli := model.Poli{ID: 2, Nama: "Poli Gigi"}
				mockRepo.On("FindByName", "Poli Gigi").Return(existingPoli, nil).Once()
			},
			expectedError: ErrPoliConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			mockRepo := new(MockPoliRepository)
			tc.setupMock(mockRepo)
			poliService := NewPoliService(mockRepo)

			_, err := poliService.UpdatePoli(context.Background(), tc.inputID, tc.inputDTO)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError))
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPoliService_DeletePoli(t *testing.T) {
	testCases := []struct {
		name          string
		inputID       int
		setupMock     func(mockRepo *MockPoliRepository)
		expectedError error
	}{
		{
			name:    "Success: Poli deleted successfully",
			inputID: 1,
			setupMock: func(mockRepo *MockPoliRepository) {
				mockRepo.On("Delete", 1).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:    "Fail: Poli not found",
			inputID: 99,
			setupMock: func(mockRepo *MockPoliRepository) {
				mockRepo.On("Delete", 99).Return(repository.ErrNotFound).Once()
			},
			expectedError: repository.ErrNotFound,
		},
		{
			name:    "Fail: Generic database error",
			inputID: 1,
			setupMock: func(mockRepo *MockPoliRepository) {
				mockRepo.On("Delete", 1).Return(errors.New("unexpected database error")).Once()
			},
			expectedError: errors.New("unexpected database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockPoliRepository)
			tc.setupMock(mockRepo)
			poliService := NewPoliService(mockRepo)

			err := poliService.DeletePoli(context.Background(), tc.inputID)

			if tc.expectedError != nil {
				assert.Error(t, err)

				if errors.Is(tc.expectedError, repository.ErrNotFound) {

					assert.True(t, errors.Is(err, tc.expectedError))
				} else {

					assert.Contains(t, err.Error(), tc.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
