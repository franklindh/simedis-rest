package model

import (
	"database/sql"
	"time"
)

type Petugas struct {
	ID        int           `json:"id,omitempty" db:"id_petugas"`
	PoliID    sql.NullInt64 `json:"poli_id" db:"id_poli"`
	Username  string        `json:"username" db:"username_petugas"`
	Name      string        `json:"name" db:"nama_petugas"`
	Status    string        `json:"status" db:"status"`
	Role      string        `json:"role" db:"role"`
	Password  string        `json:"password,omitempty" db:"password"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"update_at"`
	DeletedAt sql.NullTime  `json:"deleted_at" db:"deleted_at"`
}
