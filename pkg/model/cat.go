package model

import (
	"errors"
	"net/http"
)

type CatRequest struct {
	Name   *string `json:"cat_name"`
	Age    *int    `json:"cat_age"`
	Breed  *string `json:"cat_breed"`
	Weight *int    `json:"cat_weight"`
}

// Allow to check requested value in the body
func (a *CatRequest) Bind(r *http.Request) error {

	if a.Name == nil || *a.Name == "" {
		return errors.New("cat_name is empty")
	}

	if a.Age == nil || *a.Age <= 0 {
		return errors.New("cat_age must be a positive integer")
	}

	if a.Breed == nil || *a.Breed == "" {
		return errors.New("cat_breed is empty")
	}

	if a.Weight == nil || *a.Weight <= 0 {
		return errors.New("cat_weight must be a positive integer")
	}

	return nil
}

type CatResponse struct {
	Id     uint   `json:"id"`
	Name   string `json:"cat_name"`
	Age    int    `json:"cat_age"`
	Breed  string `json:"cat_breed"`
	Weight int    `json:"cat_weight"`
}

type CatHistoryResponse struct {
	Id     uint                    `json:"id"`
	Name   string                  `json:"cat_name"`
	Age    int                     `json:"cat_age"`
	Breed  string                  `json:"cat_breed"`
	Weight int                     `json:"cat_weight"`
	Visits []*VisitHistoryResponse `json:"cat_visits"`
}
