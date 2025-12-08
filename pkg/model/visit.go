package model

import (
	"errors"
	"net/http"
	"time"
)

type VisitRequest struct {
	CatId  *uint   `json:"visit_cat_id"`
	Date   *string `json:"visit_date"`
	Reason *string `json:"visit_reason"`
	Vet    *string `json:"visit_vet"`
}

// Allow to check requested value in the body
func (a *VisitRequest) Bind(r *http.Request) error {
	layout := "2006-01-02"

	if a.CatId == nil || *a.CatId <= 0 {
		return errors.New("visit_cat_id must be a positive integer")
	}

	if a.Date == nil || *a.Date == "" {
		return errors.New("visit_date is empty")
	}

	if a.Reason == nil || *a.Reason == "" {
		return errors.New("visit_reason is empty")
	}

	if a.Vet == nil || *a.Vet == "" {
		return errors.New("visit_vet is empty")
	}

	_, err := time.Parse(layout, *a.Date)
	if err != nil {
		return errors.New("visit_date wrong format, expected YYYY-MM-DD")
	}

	return nil
}

type VisitResponse struct {
	Id         uint                 `json:"id"`
	CatId      uint                 `json:"visit_cat_id"`
	Date       string               `json:"visit_date"`
	Reason     string               `json:"visit_reason"`
	Vet        string               `json:"visit_vet"`
	Treatments []*TreatmentResponse `json:"visit_treatments"`
}

type VisitHistoryResponse struct {
	Id         uint                        `json:"id"`
	Date       string                      `json:"visit_date"`
	Reason     string                      `json:"visit_reason"`
	Vet        string                      `json:"visit_vet"`
	Treatments []*TreatmentHistoryResponse `json:"visit_treatments"`
}
