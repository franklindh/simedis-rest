package repository

import (
	"fmt"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/pkg/utils"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {

	err := seedPoli(db)
	if err != nil {
		return err
	}

	err = seedPetugas(db)
	if err != nil {
		return err
	}

	fmt.Println("Seeding completed successfully.")
	return nil
}

func seedPoli(db *gorm.DB) error {

	var count int64
	db.Model(&model.Poli{}).Count(&count)
	if count > 0 {
		return nil
	}

	polis := []model.Poli{
		{Nama: "Poli Umum", Status: "aktif"},
		{Nama: "Poli Gigi", Status: "aktif"},
		{Nama: "Poli Anak", Status: "aktif"},
	}

	if err := db.Create(&polis).Error; err != nil {
		return fmt.Errorf("failed to seed poli: %w", err)
	}
	return nil
}

func seedPetugas(db *gorm.DB) error {

	var count int64
	db.Model(&model.Petugas{}).Count(&count)
	if count > 0 {
		return nil
	}

	hashedPassword, _ := utils.HashPassword("password123")

	petugas := []model.Petugas{
		{Username: "admin", Nama: "Admin Utama", Status: "aktif", Role: "Administrasi", Password: hashedPassword},
		{Username: "dokter.gigi", Nama: "Dr. Budi", Status: "aktif", Role: "Dokter", Password: hashedPassword},
	}

	if err := db.Create(&petugas).Error; err != nil {
		return fmt.Errorf("failed to seed petugas: %w", err)
	}
	return nil
}
