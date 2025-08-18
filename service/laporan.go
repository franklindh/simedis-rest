package service

import (
	"context"
	"fmt"
	"time"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
)

type LaporanService struct {
	repo *repository.LaporanRepository
}

func NewLaporanService(repo *repository.LaporanRepository) *LaporanService {
	return &LaporanService{repo: repo}
}

func (s *LaporanService) GetLaporanKunjunganPerPoli(ctx context.Context, startDate, endDate string) ([]model.LaporanKunjunganPoli, error) {

	layout := "2006-01-02"
	start, err1 := time.Parse(layout, startDate)
	end, err2 := time.Parse(layout, endDate)
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("invalid date format, please use YYYY-MM-DD")
	}
	if end.Before(start) {
		return nil, fmt.Errorf("endDate cannot be before startDate")
	}

	return s.repo.GetLaporanKunjunganPerPoli(startDate, endDate)
}

func (s *LaporanService) GetLaporanPenyakitTeratas(ctx context.Context, startDate, endDate string, limit int) ([]model.LaporanPenyakitTeratas, error) {
	layout := "2006-01-02"
	start, err1 := time.Parse(layout, startDate)
	end, err2 := time.Parse(layout, endDate)
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("invalid date format, please use YYYY-MM-DD")
	}
	if end.Before(start) {
		return nil, fmt.Errorf("endDate cannot be before startDate")
	}
	if limit <= 0 {
		limit = 10
	}

	return s.repo.GetLaporanPenyakitTeratas(startDate, endDate, limit)
}
