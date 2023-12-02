package dto

import "time"

type UserDto struct {
	Id          string    `json:"id"`
	Role        string    `json:"role"`
	Email       string    `json:"email"`
	IsConfirmed bool      `json:"isConfirmed"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}
