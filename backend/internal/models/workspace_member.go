package models

import "time"

type WorkspaceMember struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WorkspaceID string    `gorm:"type:uuid;not null"`
	UserID      string    `gorm:"type:uuid;not null"`
	Role        string    `gorm:"not null;default:member"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	User      User      `gorm:"foreignKey:UserID"`
	Workspace Workspace `gorm:"foreignKey:WorkspaceID"`
}
