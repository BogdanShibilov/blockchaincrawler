package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey; type:uuid; default:gen_random_uuid()"`
	Role        Role      `json:"role" gorm:"type:varchar(64); default:'user'"`
	Email       string    `json:"email" gorm:"unique; type:varchar(255); not null"`
	Password    string    `json:"password" gorm:"type:varchar(255); not null"`
	IsConfirmed bool      `json:"isConfirmed" gorm:"default:false"`
	CreatedAt   time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"default:now()"`
	Profile     Profile   `json:"profile" gorm:"foreignKey:UserId"`
}

type Profile struct {
	UserId  uuid.UUID `json:"userId" gorm:"type:uuid;"`
	Name    string    `json:"name" gorm:"type:varchar(64); default:'';"`
	Surname string    `json:"surname" gorm:"type:varchar(64); default:'';"`
	AboutMe string    `json:"aboutMe" gorm:"type:varchar(255); default:'';"`
}

func (p *Profile) FullName() string {
	return p.Name + p.Surname
}
