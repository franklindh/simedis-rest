package service

import (
	"time"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
)

type AntrianRepository interface {
	Create(antrian model.Antrian) (model.Antrian, error)
	GetAll(params repository.ParamsGetAllAntrian) ([]model.Antrian, pagination.Metadata, error)
	GetByID(id int) (model.Antrian, error)
	Update(id int, antrian model.Antrian) (model.Antrian, error)
	Delete(id int) error
	CheckAntrian(pasienID, jadwalID int) (bool, error)
	CheckForOverlappingAntrian(pasienID int, tanggal, waktuMulai, waktuSelesai time.Time) (bool, error)
	CountTodayByJadwal(jadwalID int) (int64, error)
}

type JadwalRepository interface {
	Create(jadwal model.Jadwal) (model.Jadwal, error)
	GetAll(params repository.ParamsGetAllJadwal) ([]model.Jadwal, pagination.Metadata, error)
	GetById(id int) (model.Jadwal, error)
	Update(id int, jadwal model.Jadwal) (model.Jadwal, error)
	Delete(id int) error
}

type JenisPemeriksaanLabRepository interface {
	Create(jenis model.JenisPemeriksaanLab) (model.JenisPemeriksaanLab, error)
	GetAll() ([]model.JenisPemeriksaanLab, error)
	GetById(id int) (model.JenisPemeriksaanLab, error)
	Update(id int, jenis model.JenisPemeriksaanLab) (model.JenisPemeriksaanLab, error)
	Delete(id int) error
	FindByName(name string) (model.JenisPemeriksaanLab, error)
}

type PasienRepository interface {
	Create(pasien model.Pasien) (model.Pasien, error)
	GetAll(params repository.ParamsGetAllPasien) ([]model.Pasien, pagination.Metadata, error)
	GetById(id int) (model.Pasien, error)
	Update(id int, pasien model.Pasien) (model.Pasien, error)
	Delete(id int) error
	GetLastID() (int, error)
}

type PemeriksaanLabRepository interface {
	Create(hasilLab model.PemeriksaanLab) (model.PemeriksaanLab, error)
	GetAllByPemeriksaanID(pemeriksaanID int) ([]model.PemeriksaanLab, error)
	GetById(id int) (model.PemeriksaanLab, error)
	Update(id int, hasilLab model.PemeriksaanLab) (model.PemeriksaanLab, error)
	Delete(id int) error
}

type PemeriksaanRepository interface {
	CheckExistingPemeriksaan(antrianID int) error
	Create(pemeriksaan model.Pemeriksaan) (model.Pemeriksaan, error)
	GetById(id int) (model.Pemeriksaan, error)
	GetAllByPasienID(pasienID int) ([]model.Pemeriksaan, error)
	Update(id int, pemeriksaan model.Pemeriksaan) (model.Pemeriksaan, error)
	Delete(id int) error
}

type PetugasRepository interface {
	GetByUsername(username string) (model.Petugas, error)
	Create(petugas model.Petugas) (model.Petugas, error)
	GetAll(params repository.ParamsGetAllPetugas) ([]model.Petugas, pagination.Metadata, error)
	GetById(id int) (model.Petugas, error)
	Update(id int, petugas model.Petugas) (model.Petugas, error)
	Delete(id int) error
	UpdatePassword(id int, newHashedPassword string) error
}

type PoliRepository interface {
	Create(poli model.Poli) (model.Poli, error)
	GetAll() ([]model.Poli, error)
	GetById(id int) (model.Poli, error)
	Update(id int, poli model.Poli) (model.Poli, error)
	Delete(id int) error
	FindByName(name string) (model.Poli, error)
}
