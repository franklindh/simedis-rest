package service

import (
	"context"
	"fmt"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/jackc/pgx/v5/pgconn"
)

type PemeriksaanLabService struct {
	repo PemeriksaanLabRepository
}

func NewPemeriksaanLabService(repo PemeriksaanLabRepository) *PemeriksaanLabService {
	return &PemeriksaanLabService{repo: repo}
}

func (s *PemeriksaanLabService) CreateBatch(ctx context.Context, pemeriksaanID int, reqs []model.CreateHasilLabRequest) ([]model.PemeriksaanLab, error) {
	var createdResults []model.PemeriksaanLab
	for _, req := range reqs {
		hasilLab := model.PemeriksaanLab{
			PemeriksaanID:      pemeriksaanID,
			JenisPemeriksaanID: req.JenisPemeriksaanID,
			Hasil:              req.Hasil,
		}
		created, err := s.repo.Create(hasilLab)
		if err != nil {

			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23503" {
				return nil, fmt.Errorf("invalid data: %s", pgErr.Detail)
			}
			return nil, fmt.Errorf("failed to create lab result for jenis_id %d: %w", req.JenisPemeriksaanID, err)
		}
		createdResults = append(createdResults, created)
	}
	return createdResults, nil
}

func (s *PemeriksaanLabService) GetAllByPemeriksaanID(ctx context.Context, pemeriksaanID int) ([]model.PemeriksaanLab, error) {
	return s.repo.GetAllByPemeriksaanID(pemeriksaanID)
}

func (s *PemeriksaanLabService) Update(ctx context.Context, id int, req model.UpdateHasilLabRequest) (model.PemeriksaanLab, error) {
	hasilLab := model.PemeriksaanLab{
		Hasil: req.Hasil,
	}
	return s.repo.Update(id, hasilLab)
}

func (s *PemeriksaanLabService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
