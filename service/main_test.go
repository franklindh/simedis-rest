package service

import (
	"time"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/internal/repository"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"github.com/stretchr/testify/mock"
)

type MockAntrianRepository struct {
	mock.Mock
}

var _ AntrianRepository = (*MockAntrianRepository)(nil)

func (m *MockAntrianRepository) Create(antrian model.Antrian) (model.Antrian, error) {
	args := m.Called(antrian)
	return args.Get(0).(model.Antrian), args.Error(1)
}
func (m *MockAntrianRepository) GetAll(params repository.ParamsGetAllAntrian) ([]model.Antrian, pagination.Metadata, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(pagination.Metadata), args.Error(2)
	}
	return args.Get(0).([]model.Antrian), args.Get(1).(pagination.Metadata), args.Error(2)
}
func (m *MockAntrianRepository) GetByID(id int) (model.Antrian, error) {
	args := m.Called(id)
	return args.Get(0).(model.Antrian), args.Error(1)
}
func (m *MockAntrianRepository) Update(id int, antrian model.Antrian) (model.Antrian, error) {
	args := m.Called(id, antrian)
	return args.Get(0).(model.Antrian), args.Error(1)
}
func (m *MockAntrianRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockAntrianRepository) CheckAntrian(pasienID, jadwalID int) (bool, error) {
	args := m.Called(pasienID, jadwalID)
	return args.Bool(0), args.Error(1)
}
func (m *MockAntrianRepository) CheckForOverlappingAntrian(pasienID int, tanggal time.Time, waktuMulai time.Time, waktuSelesai time.Time) (bool, error) {
	args := m.Called(pasienID, tanggal, waktuMulai, waktuSelesai)
	return args.Bool(0), args.Error(1)
}
func (m *MockAntrianRepository) CountTodayByJadwal(jadwalID int) (int64, error) {
	args := m.Called(jadwalID)
	return args.Get(0).(int64), args.Error(1)
}

type MockJadwalRepository struct {
	mock.Mock
}

var _ JadwalRepository = (*MockJadwalRepository)(nil)

func (m *MockJadwalRepository) Create(jadwal model.Jadwal) (model.Jadwal, error) {
	args := m.Called(jadwal)
	return args.Get(0).(model.Jadwal), args.Error(1)
}
func (m *MockJadwalRepository) GetAll(params repository.ParamsGetAllJadwal) ([]model.Jadwal, pagination.Metadata, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(pagination.Metadata), args.Error(2)
	}
	return args.Get(0).([]model.Jadwal), args.Get(1).(pagination.Metadata), args.Error(2)
}
func (m *MockJadwalRepository) GetById(id int) (model.Jadwal, error) {
	args := m.Called(id)
	return args.Get(0).(model.Jadwal), args.Error(1)
}
func (m *MockJadwalRepository) Update(id int, jadwal model.Jadwal) (model.Jadwal, error) {
	args := m.Called(id, jadwal)
	return args.Get(0).(model.Jadwal), args.Error(1)
}
func (m *MockJadwalRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

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

type MockPemeriksaanRepository struct {
	mock.Mock
}

var _ PemeriksaanRepository = (*MockPemeriksaanRepository)(nil)

func (m *MockPemeriksaanRepository) CheckExistingPemeriksaan(antrianID int) error {
	args := m.Called(antrianID)
	return args.Error(0)
}
func (m *MockPemeriksaanRepository) Create(pemeriksaan model.Pemeriksaan) (model.Pemeriksaan, error) {
	args := m.Called(pemeriksaan)
	return args.Get(0).(model.Pemeriksaan), args.Error(1)
}
func (m *MockPemeriksaanRepository) GetById(id int) (model.Pemeriksaan, error) {
	args := m.Called(id)
	return args.Get(0).(model.Pemeriksaan), args.Error(1)
}
func (m *MockPemeriksaanRepository) GetAllByPasienID(pasienID int) ([]model.Pemeriksaan, error) {
	args := m.Called(pasienID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Pemeriksaan), args.Error(1)
}
func (m *MockPemeriksaanRepository) Update(id int, pemeriksaan model.Pemeriksaan) (model.Pemeriksaan, error) {
	args := m.Called(id, pemeriksaan)
	return args.Get(0).(model.Pemeriksaan), args.Error(1)
}
func (m *MockPemeriksaanRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockPetugasRepository struct {
	mock.Mock
}

func (m *MockPetugasRepository) Create(petugas model.Petugas) (model.Petugas, error) {
	args := m.Called(petugas)
	return args.Get(0).(model.Petugas), args.Error(1)
}

func (m *MockPetugasRepository) GetByUsername(username string) (model.Petugas, error) {
	args := m.Called(username)
	return args.Get(0).(model.Petugas), args.Error(1)
}

func (m *MockPetugasRepository) GetById(id int) (model.Petugas, error) {
	args := m.Called(id)
	return args.Get(0).(model.Petugas), args.Error(1)
}

func (m *MockPetugasRepository) GetAll(params repository.ParamsGetAllPetugas) ([]model.Petugas, pagination.Metadata, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Get(1).(pagination.Metadata), args.Error(2)
	}
	return args.Get(0).([]model.Petugas), args.Get(1).(pagination.Metadata), args.Error(2)
}

func (m *MockPetugasRepository) Update(id int, petugas model.Petugas) (model.Petugas, error) {
	args := m.Called(id, petugas)
	return args.Get(0).(model.Petugas), args.Error(1)
}

func (m *MockPetugasRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockPoliRepository struct {
	mock.Mock
}

var _ PoliRepository = (*MockPoliRepository)(nil)

func (m *MockPoliRepository) Create(poli model.Poli) (model.Poli, error) {
	args := m.Called(poli)
	return args.Get(0).(model.Poli), args.Error(1)
}

func (m *MockPoliRepository) GetAll() ([]model.Poli, error) {
	args := m.Called()
	return args.Get(0).([]model.Poli), args.Error(1)
}

func (m *MockPoliRepository) GetById(id int) (model.Poli, error) {
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
