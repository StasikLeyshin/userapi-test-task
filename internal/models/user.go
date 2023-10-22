package models

import (
	"errors"
	"time"
)

var (
	UserNotFound = errors.New("user_not_found")
	EmptyValues  = errors.New("empty_values")
)

type User struct {
	ID          int       `json:"user_id"`
	CreateDate  time.Time `json:"create_date"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
}

type UserList map[string]User

type UserStore struct {
	Increment int      `json:"increment"`
	List      UserList `json:"list"`
}

type CreateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}
