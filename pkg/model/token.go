package model

import (
	"errors"
	"net/http"
)

type RefreshTokenRequest struct {
	RefreshToken *string `json:"refresh_token"`
}

// Allow to check requested value in the body
func (a *RefreshTokenRequest) Bind(r *http.Request) error {

	if a.RefreshToken == nil || *a.RefreshToken == "" {
		return errors.New("refresh_token is empty")
	}

	return nil
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
