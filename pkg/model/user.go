package model

import (
	"errors"
	"net/http"
)

type UserRequest struct {
	Email    *string `json:"user_email"`
	Password *string `json:"user_password"`
	Role     *string `json:"user_role"`
}

type UserLoginRequest struct {
	Email    *string `json:"user_email"`
	Password *string `json:"user_password"`
}

// Allow to check requested value in the body
func (a *UserRequest) Bind(r *http.Request) error {

	if a.Email == nil || *a.Email == "" {
		return errors.New("user_email is empty")
	}

	if a.Password == nil || *a.Password == "" {
		return errors.New("user_password is empty")
	}

	if a.Role == nil || *a.Role == "" {
		return errors.New("user_role is empty")
	}

	return nil
}

// Allow to check requested value in the body
func (a *UserLoginRequest) Bind(r *http.Request) error {

	if a.Email == nil || *a.Email == "" {
		return errors.New("user_email is empty")
	}

	if a.Password == nil || *a.Password == "" {
		return errors.New("user_password is empty")
	}

	return nil
}

type UserResponse struct {
	Id    uint   `json:"id"`
	Email string `json:"user_email"`
	Role  string `json:"user_role"`
}
