package repository

import (
	"database/sql"
	"errors"

	"github.com/franklindh/simedis-api/internal/model"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("data not found")

type PoliRepository struct {
	DB *gorm.DB
}

func NewPoliRepository(db *gorm.DB) *PoliRepository {
	return &PoliRepository{DB: db}
}

func (r *PoliRepository) GetAll() ([]model.Poli, error) {
	// pake sql manual
	// query := `SELECT id_poli, nama_poli, status_poli, created_at, updated_at FROM poli WHERE deleted_at IS NULL`

	// rows, err := r.DB.Query(query)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()

	// var poli []model.Poli
	// for rows.Next() {
	// 	var p model.Poli
	// 	if err := rows.Scan(&p.ID, &p.Name, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
	// 		return nil, err
	// 	}
	// 	poli = append(poli, p)
	// }
	// return poli, nil

	// pake gorm
	var poli []model.Poli
	result := r.DB.Find(&poli)
	return poli, result.Error
}

func (r *PoliRepository) GetByID(id int) (model.Poli, error) {
	// pake sql manual
	// query := `SELECT id_poli, nama_poli, status_poli, created_at, updated_at FROM poli WHERE id_poli = $1 AND deleted_at IS NULL`

	// var p model.Poli
	// err := r.DB.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return model.Poli{}, ErrNotFound
	// 	}
	// 	return model.Poli{}, err
	// }
	// return p, nil

	// pake gorm
	var poli model.Poli
	result := r.DB.First(&poli, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Poli{}, ErrNotFound
		}
		return model.Poli{}, result.Error
	}
	return poli, nil
}

func (r *PoliRepository) Create(poli model.Poli) (model.Poli, error) {
	// pake sql manual
	// query := `INSERT INTO poli (nama_poli, status_poli) VALUES ($1, $2) RETURNING id_poli, created_at, updated_at`

	// err := r.DB.QueryRow(query, poli.Name, poli.Status).Scan(&poli.ID, &poli.CreatedAt, &poli.UpdatedAt)
	// if err != nil {
	// 	return model.Poli{}, err
	// }

	// return poli, nil

	// pake gorm
	result := r.DB.Create(&poli)
	return poli, result.Error
}

func (r *PoliRepository) Update(id int, poli model.Poli) (model.Poli, error) {
	// pake sql manual
	// query := `UPDATE poli SET nama_poli = $1, status_poli = $2, updated_at = NOW() WHERE id_poli = $3 RETURNING id_poli, nama_poli, status_poli, created_at, updated_at`

	// var updatePoli model.Poli

	// err := r.DB.QueryRow(query, poli.Name, poli.Status, id).Scan(&updatePoli.ID, &updatePoli.Name, &updatePoli.Status, &updatePoli.CreatedAt, &updatePoli.UpdatedAt)

	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return model.Poli{}, ErrNotFound
	// 	}
	// 	return model.Poli{}, err
	// }

	// return updatePoli, nil

	// pake gorm
	if _, err := r.GetByID(id); err != nil {
		return model.Poli{}, err
	}
	poli.ID = id
	result := r.DB.Save(&poli)
	return poli, result.Error
}

func (r *PoliRepository) Delete(id int) error {
	// pake sql manual
	// query := `UPDATE poli SET deleted_at = NOW() WHERE id_poli = $1 AND deleted_at IS NULL`

	// res, err := r.DB.Exec(query, id)
	// if err != nil {
	// 	return err
	// }

	// rowsAffedcted, err := res.RowsAffected()
	// if err != nil {
	// 	return err
	// }

	// if rowsAffedcted == 0 {
	// 	return ErrNotFound
	// }

	// return nil

	// pake gorm
	result := r.DB.Delete(&model.Poli{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PoliRepository) FindByNameIncludingDeleted(name string) (model.Poli, error) {
	// pake sql manual
	// query := `SELECT id_poli, nama_poli, status_poli, created_at, updated_at, deleted_at FROM poli WHERE nama_poli = $1`

	// var p model.Poli
	// err := r.DB.QueryRow(query, name).Scan(&p.ID, &p.Name, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	// if err != nil {
	// 	return model.Poli{}, err
	// }
	// return p, nil

	// pake gorm
	var poli model.Poli

	result := r.DB.Unscoped().Where("nama_poli = ?", name).First(&poli)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Poli{}, sql.ErrNoRows
		}
		return model.Poli{}, result.Error
	}
	return poli, nil
}

func (r *PoliRepository) Restore(id int) (model.Poli, error) {
	// pake sql manual
	// query := `UPDATE poli SET deleted_at = NULL, updated_at = NOW() WHERE id_poli = $1 AND deleted_at IS NOT NULL RETURNING id_poli, nama_poli, status_poli, created_at, updated_at`

	// var restoredPoli model.Poli
	// err := r.DB.QueryRow(query, id).Scan(&restoredPoli.ID, &restoredPoli.Name, &restoredPoli.Status, &restoredPoli.CreatedAt, &restoredPoli.UpdatedAt)

	// return restoredPoli, err

	// pake gorm
	var poli model.Poli

	result := r.DB.Unscoped().First(&poli, id)
	if result.Error != nil {
		return model.Poli{}, ErrNotFound
	}
	return poli, nil
}
