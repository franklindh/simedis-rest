package model

import (
	"database/sql"
	"time"
)

type Poli struct {
	ID        int          `json:"id,omitempty" db:"id_poli"`
	Name      string       `json:"nama" db:"nama_poli"`
	Status    string       `json:"status" db:"status_poli"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"delete_at" db:"deleted_at"`
}
