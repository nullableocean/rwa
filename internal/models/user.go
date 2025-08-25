package models

import (
	"errors"
	"time"
)

type User struct {
	ID             int64     `json:"-"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Bio            string    `json:"bio"`
	Image          string    `json:"image"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	HashedPassword string    `json:"-"`
}

type UserUpdateInfo struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type UserCreateInfo struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Password string `json:"password"`
}

func (info *UserCreateInfo) Validate() error {
	switch "" {
	case info.Email:
		return errors.New("Email is empty")
	case info.Username:
		return errors.New("Username is empty")
	case info.Password:
		return errors.New("Password is empty")
	}

	return nil
}
