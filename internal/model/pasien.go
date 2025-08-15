package model

import (
	"database/sql"
	"time"
)

type Pasien struct {
	ID                        int            `json:"id,omitempty" gorm:"primaryKey;column:id_pasien"`
	NIK                       string         `json:"nik" gorm:"column:nik;unique" binding:"required,numeric,min=16,max=16"`
	NoRekamMedis              sql.NullString `json:"no_rekam_medis,omitempty" gorm:"column:no_rekam_medis;unique"`
	NoKartuJaminan            sql.NullString `json:"no_kartu_jaminan,omitempty" gorm:"column:no_kartu_jaminan"`
	UsernamePasien            string         `json:"username_pasien" gorm:"column:username_pasien;unique" binding:"required,alphanum,min=5"`
	NoTeleponPasien           sql.NullString `json:"no_telepon_pasien,omitempty" gorm:"column:no_telepon_pasien"`
	NamaPasien                string         `json:"nama_pasien" gorm:"column:nama_pasien" binding:"required"`
	AlamatPasien              string         `json:"alamat_pasien" gorm:"column:alamat_pasien" binding:"required"`
	TempatLahirPasien         string         `json:"tempat_lahir_pasien" gorm:"column:tempat_lahir_pasien" binding:"required"`
	TanggalLahirPasien        string         `json:"tanggal_lahir_pasien" gorm:"column:tanggal_lahir_pasien" binding:"required,datetime=2006-01-02"`
	JKPasien                  string         `json:"jk_pasien" gorm:"column:jk_pasien" binding:"required,oneof=L P"`
	StatusPernikahan          string         `json:"status_pernikahan" gorm:"column:status_pernikahan" binding:"required,oneof='Belum Menikah' Menikah 'Cerai Hidup' 'Cerai Mati'"`
	NamaKeluargaTerdekat      sql.NullString `json:"nama_keluarga_terdekat,omitempty" gorm:"column:nama_keluarga_terdekat"`
	NoTeleponKeluargaTerdekat sql.NullString `json:"no_telepon_keluarga_terdekat,omitempty" gorm:"column:no_telepon_keluarga_terdekat"`
	Password                  string         `json:"-" gorm:"column:password"`
	CreatedAt                 time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt                 time.Time      `json:"updated_at" gorm:"column:updated_at"`
}

func (Pasien) TableName() string {
	return "pasien"
}
