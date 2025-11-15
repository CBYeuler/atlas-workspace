package models

import "time"

type Workspace struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerID   string    `gorm:"type:uuid;not null"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Owner User `gorm:"foreignKey:OwnerID"`
}
