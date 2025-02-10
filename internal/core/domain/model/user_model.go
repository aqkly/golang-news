package model

import "time"

type User struct {
	ID        int64      `gorm:"id"`
	Nama      string     `gorm:"nama"`
	Email     string     `gorm:"email"`
	Password  string     `gorm:"password"`
	CreatedAt time.Time  `gorm:"created_at"`
	UpdatedAt *time.Time `gorm:"updated_at"`
}
