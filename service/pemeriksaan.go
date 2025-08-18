package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
)

var (
	ErrPemeriksaanExists = errors.New("pemeriksaan for this antrian already exists")
)

type PemeriksaanService struct {
	repo        *repository.PemeriksaanRepository
	antrianRepo *repository.AntrianRepository
}

func NewPemeriksaanService(repo *repository.PemeriksaanRepository, antrianRepo *repository.AntrianRepository) *PemeriksaanService {
	return &PemeriksaanService{repo: repo, antrianRepo: antrianRepo}
}

func (s *PemeriksaanService) CreatePemeriksaan(ctx context.Context, req model.CreatePemeriksaanRequest) (model.Pemeriksaan, error) {

	err := s.repo.CheckExistingPemeriksaan(req.AntrianID)
	if err == nil {
		return model.Pemeriksaan{}, ErrPemeriksaanExists
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return model.Pemeriksaan{}, fmt.Errorf("error checking existing pemeriksaan: %w", err)
	}

	pemeriksaan := model.Pemeriksaan{
		AntrianID:          req.AntrianID,
		Nadi:               sql.NullString{String: req.Nadi, Valid: req.Nadi != ""},
		TekananDarah:       sql.NullString{String: req.TekananDarah, Valid: req.TekananDarah != ""},
		Suhu:               sql.NullString{String: req.Suhu, Valid: req.Suhu != ""},
		BeratBadan:         sql.NullString{String: req.BeratBadan, Valid: req.BeratBadan != ""},
		KeadaanUmum:        sql.NullString{String: req.KeadaanUmum, Valid: req.KeadaanUmum != ""},
		Keluhan:            sql.NullString{String: req.Keluhan, Valid: req.Keluhan != ""},
		RiwayatPenyakit:    sql.NullString{String: req.RiwayatPenyakit, Valid: req.RiwayatPenyakit != ""},
		Keterangan:         sql.NullString{String: req.Keterangan, Valid: req.Keterangan != ""},
		Tindakan:           sql.NullString{String: req.Tindakan, Valid: req.Tindakan != ""},
		TanggalPemeriksaan: req.TanggalPemeriksaan,
	}
	if req.IcdID != nil {
		pemeriksaan.IcdID = sql.NullInt64{Int64: int64(*req.IcdID), Valid: true}
	}

	createdPemeriksaan, err := s.repo.Create(pemeriksaan)
	if err != nil {
		return model.Pemeriksaan{}, err
	}

	antrian, err := s.antrianRepo.GetByID(req.AntrianID)
	if err == nil {
		antrian.Status = "Selesai"
		s.antrianRepo.Update(antrian.ID, antrian)
	}

	return createdPemeriksaan, nil
}

func (s *PemeriksaanService) GetPemeriksaanByID(ctx context.Context, id int) (model.Pemeriksaan, error) {
	return s.repo.GetByID(id)
}

func (s *PemeriksaanService) GetRiwayatPemeriksaanPasien(ctx context.Context, pasienID int) ([]model.PemeriksaanDetail, error) {
	allPemeriksaan, err := s.repo.GetAllByPasienID(pasienID)
	if err != nil {
		return nil, err
	}

	var responseData []model.PemeriksaanDetail
	for _, p := range allPemeriksaan {
		detail := model.PemeriksaanDetail{
			ID:                 p.ID,
			TanggalPemeriksaan: p.TanggalPemeriksaan,
			Nadi:               p.Nadi,
			TekananDarah:       p.TekananDarah,
			Suhu:               p.Suhu,
			BeratBadan:         p.BeratBadan,
			Keluhan:            p.Keluhan,
			Tindakan:           p.Tindakan,
			Pasien: struct {
				ID           int    `json:"id"`
				Nama         string `json:"nama"`
				NoRekamMedis string `json:"no_rekam_medis"`
			}{
				ID:           p.Antrian.Pasien.ID,
				Nama:         p.Antrian.Pasien.NamaPasien,
				NoRekamMedis: p.Antrian.Pasien.NoRekamMedis.String,
			},
			Dokter: struct {
				ID   int    `json:"id"`
				Nama string `json:"nama"`
			}{
				ID:   p.Antrian.Jadwal.Petugas.ID,
				Nama: p.Antrian.Jadwal.Petugas.Nama,
			},
			Poli: struct {
				ID   int    `json:"id"`
				Nama string `json:"nama"`
			}{
				ID:   p.Antrian.Jadwal.Poli.ID,
				Nama: p.Antrian.Jadwal.Poli.Nama,
			},
			Diagnosis: struct {
				ID       sql.NullInt64  `json:"id"`
				Kode     sql.NullString `json:"kode"`
				Penyakit sql.NullString `json:"penyakit"`
			}{

				ID: sql.NullInt64{Int64: int64(p.Icd.ID), Valid: p.IcdID.Valid},

				Kode:     sql.NullString{String: p.Icd.KodeIcd, Valid: p.IcdID.Valid},
				Penyakit: sql.NullString{String: p.Icd.NamaPenyakit, Valid: p.IcdID.Valid},
			},
		}
		responseData = append(responseData, detail)
	}
	return responseData, nil
}

func (s *PemeriksaanService) UpdatePemeriksaan(ctx context.Context, id int, req model.UpdatePemeriksaanRequest) (model.Pemeriksaan, error) {
	pemeriksaanUpdate := model.Pemeriksaan{
		Nadi:               sql.NullString{String: req.Nadi, Valid: req.Nadi != ""},
		TekananDarah:       sql.NullString{String: req.TekananDarah, Valid: req.TekananDarah != ""},
		Suhu:               sql.NullString{String: req.Suhu, Valid: req.Suhu != ""},
		BeratBadan:         sql.NullString{String: req.BeratBadan, Valid: req.BeratBadan != ""},
		KeadaanUmum:        sql.NullString{String: req.KeadaanUmum, Valid: req.KeadaanUmum != ""},
		Keluhan:            sql.NullString{String: req.Keluhan, Valid: req.Keluhan != ""},
		RiwayatPenyakit:    sql.NullString{String: req.RiwayatPenyakit, Valid: req.RiwayatPenyakit != ""},
		Keterangan:         sql.NullString{String: req.Keterangan, Valid: req.Keterangan != ""},
		Tindakan:           sql.NullString{String: req.Tindakan, Valid: req.Tindakan != ""},
		TanggalPemeriksaan: req.TanggalPemeriksaan,
	}
	if req.IcdID != nil {
		pemeriksaanUpdate.IcdID = sql.NullInt64{Int64: int64(*req.IcdID), Valid: true}
	}

	return s.repo.Update(id, pemeriksaanUpdate)
}

func (s *PemeriksaanService) DeletePemeriksaan(ctx context.Context, id int) error {
	return s.repo.Delete(id)
}
