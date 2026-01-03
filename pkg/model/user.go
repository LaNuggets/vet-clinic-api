package model

import (
	"errors"
	"net/http"
)

type UserRequest struct {
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

	return nil
}

type UserResponse struct {
	Id       uint   `json:"id"`
	Email    string `json:"user_email"`
	Password string `json:"user_password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
