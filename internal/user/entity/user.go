package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey; type:uuid; default:gen_random_uuid()"`
	Email       string    `json:"email" gorm:"unique; type:varchar(255); not null"`
	Password    string    `json:"password" gorm:"type:varchar(255); not null"`
	IsConfirmed bool      `json:"isConfirmed" gorm:"default:false"`
	CreatedAt   time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"default:now()"`
}
