package service

import (
	"context"
	"errors"
	"testing"

	"github.com/franklindh/simedis-api/internal/config"
	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPetugasService_CreatePetugas(t *testing.T) {
	cfg := &config.Config{
		DefaultPetugasPassword: "passworddefault",
	}

	mockRepo := new(MockPetugasRepository)
	petugasService := NewPetugasService(mockRepo, cfg)

	inputDTO := model.CreatePetugasRequest{
		Username: "johndoe",
		Nama:     "John Doe",
		Role:     "Dokter",
		Status:   "aktif",
	}

	t.Run("Success: Petugas created successfully", func(t *testing.T) {

		mockReturnPetugas := model.Petugas{
			ID:       1,
			Username: "johndoe",
			Nama:     "John Doe",
			Role:     "Dokter",
			Status:   "aktif",
		}

		mockRepo.On("Create", mock.AnythingOfType("model.Petugas")).Return(mockReturnPetugas, nil).Once()

		result, err := petugasService.CreatePetugas(context.Background(), inputDTO)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		assert.Equal(t, "johndoe", result.Username)
		assert.Equal(t, 1, result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Username conflict", func(t *testing.T) {
		pgErr := &pgconn.PgError{Code: "23505"}
		mockRepo.On("Create", mock.AnythingOfType("model.Petugas")).Return(model.Petugas{}, pgErr).Once()

		result, err := petugasService.CreatePetugas(context.Background(), inputDTO)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrPetugasConflict))
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestPetugasService_Login(t *testing.T) {
	hashedDefaultPassword, _ := utils.HashPassword("passworddefault")
	cfg := &config.Config{
		JWTSecret: "secretkeyrahasia",
	}

	testCases := []struct {
		name          string
		input         model.LoginPetugasRequest
		setupMock     func(mockRepo *MockPetugasRepository)
		expectedError error
	}{
		{
			name:  "Success: Login successful",
			input: model.LoginPetugasRequest{Username: "johndoe", Password: "passworddefault"},
			setupMock: func(mockRepo *MockPetugasRepository) {
				userFromDB := model.Petugas{
					ID:       1,
					Username: "johndoe",
					Password: hashedDefaultPassword,
					Role:     "Dokter",
				}
				mockRepo.On("GetByUsername", "johndoe").Return(userFromDB, nil).Once()
			},
			expectedError: nil,
		},
		{
			name:  "Fail: User not found",
			input: model.LoginPetugasRequest{Username: "notfound", Password: "passworddefault"},
			setupMock: func(mockRepo *MockPetugasRepository) {
				mockRepo.On("GetByUsername", "notfound").Return(model.Petugas{}, repository.ErrNotFound).Once()
			},
			expectedError: ErrInvalidCredentials,
		},
		{
			name:  "Fail: Wrong password",
			input: model.LoginPetugasRequest{Username: "johndoe", Password: "wrongpassword"},
			setupMock: func(mockRepo *MockPetugasRepository) {
				userFromDB := model.Petugas{
					ID:       1,
					Username: "johndoe",
					Password: hashedDefaultPassword,
				}
				mockRepo.On("GetByUsername", "johndoe").Return(userFromDB, nil).Once()
			},
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockPetugasRepository)
			tc.setupMock(mockRepo)
			petugasService := NewPetugasService(mockRepo, cfg)

			token, err := petugasService.Login(context.Background(), tc.input)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.expectedError))
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPetugasService_GetPetugasByID(t *testing.T) {
	mockRepo := new(MockPetugasRepository)
	petugasService := NewPetugasService(mockRepo, &config.Config{})

	t.Run("Success: Petugas found", func(t *testing.T) {
		mockPetugas := model.Petugas{ID: 1, Username: "testuser", Nama: "Test User"}
		mockRepo.On("GetById", 1).Return(mockPetugas, nil).Once()

		result, err := petugasService.GetPetugasByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.ID)
		assert.Equal(t, "testuser", result.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Petugas not found", func(t *testing.T) {
		mockRepo.On("GetById", 2).Return(model.Petugas{}, repository.ErrNotFound).Once()

		result, err := petugasService.GetPetugasByID(context.Background(), 2)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestPetugasService_UpdatePetugas(t *testing.T) {
	mockRepo := new(MockPetugasRepository)
	petugasService := NewPetugasService(mockRepo, &config.Config{})

	inputDTO := model.UpdatePetugasRequest{
		Nama:   "John Doe Updated",
		Status: "nonaktif",
		Role:   "Admin",
	}

	t.Run("Success: Petugas updated", func(t *testing.T) {
		mockReturnPetugas := model.Petugas{
			ID:     1,
			Nama:   "John Doe Updated",
			Status: "nonaktif",
			Role:   "Admin",
		}

		mockRepo.On("Update", 1, mock.AnythingOfType("model.Petugas")).Return(mockReturnPetugas, nil).Once()

		result, err := petugasService.UpdatePetugas(context.Background(), 1, inputDTO)

		assert.NoError(t, err)
		assert.Equal(t, "John Doe Updated", result.Nama)
		assert.Equal(t, "nonaktif", result.Status)
		mockRepo.AssertExpectations(t)
	})
}

func TestPetugasService_GetAllPetugas(t *testing.T) {
	mockRepo := new(MockPetugasRepository)
	petugasService := NewPetugasService(mockRepo, &config.Config{})
	params := repository.ParamsGetAllPetugas{Page: 1, PageSize: 5}

	t.Run("Success: Get all petugas", func(t *testing.T) {
		mockPetugasList := []model.Petugas{
			{ID: 1, Nama: "User Satu"},
			{ID: 2, Nama: "User Dua"},
		}
		mockMetadata := pagination.Metadata{TotalRecords: 2}

		mockRepo.On("GetAll", params).Return(mockPetugasList, mockMetadata, nil).Once()

		result, _, err := petugasService.GetAllPetugas(context.Background(), params)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "User Satu", result[0].Nama)
		mockRepo.AssertExpectations(t)
	})
}

func TestPetugasService_DeletePetugas(t *testing.T) {
	mockRepo := new(MockPetugasRepository)
	petugasService := NewPetugasService(mockRepo, &config.Config{})

	t.Run("Success: Petugas deleted", func(t *testing.T) {
		mockRepo.On("Delete", 1).Return(nil).Once()
		err := petugasService.DeletePetugas(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Petugas not found", func(t *testing.T) {
		mockRepo.On("Delete", 2).Return(repository.ErrNotFound).Once()
		err := petugasService.DeletePetugas(context.Background(), 2)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}

func TestPetugasService_ChangePassword(t *testing.T) {

	petugasID := 1
	oldPassword := "passwordLama123"
	hashedOldPassword, _ := utils.HashPassword(oldPassword)

	req := model.ChangePasswordRequest{
		OldPassword:     oldPassword,
		NewPassword:     "passwordBaru456",
		ConfirmPassword: "passwordBaru456",
	}

	t.Run("Success: Change password successfully", func(t *testing.T) {

		mockRepo := new(MockPetugasRepository)
		cfg := &config.Config{}
		service := NewPetugasService(mockRepo, cfg)

		mockPetugas := model.Petugas{ID: petugasID, Password: hashedOldPassword}
		mockRepo.On("GetById", petugasID).Return(mockPetugas, nil).Once()
		mockRepo.On("UpdatePassword", petugasID, mock.AnythingOfType("string")).Return(nil).Once()

		err := service.ChangePassword(context.Background(), petugasID, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: Old password does not match", func(t *testing.T) {

		mockRepo := new(MockPetugasRepository)
		cfg := &config.Config{}
		service := NewPetugasService(mockRepo, cfg)

		mockPetugas := model.Petugas{ID: petugasID, Password: "hash_yang_berbeda"}
		mockRepo.On("GetById", petugasID).Return(mockPetugas, nil).Once()

		err := service.ChangePassword(context.Background(), petugasID, req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrOldPasswordMismatch))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Fail: User not found", func(t *testing.T) {

		mockRepo := new(MockPetugasRepository)
		cfg := &config.Config{}
		service := NewPetugasService(mockRepo, cfg)

		mockRepo.On("GetById", petugasID).Return(model.Petugas{}, repository.ErrNotFound).Once()

		err := service.ChangePassword(context.Background(), petugasID, req)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrNotFound))
		mockRepo.AssertExpectations(t)
	})
}
