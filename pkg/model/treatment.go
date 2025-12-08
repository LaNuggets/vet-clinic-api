package model

import (
	"errors"
	"net/http"
)

type TreatmentRequest struct {
	Name    *string `json:"treatment_name"`
	VisitId *uint   `json:"treatment_visit_id"`
}

// Allow to check requested value in the body
func (a *TreatmentRequest) Bind(r *http.Request) error {

	if a.Name == nil || *a.Name == "" {
		return errors.New("treatment_name is empty")
	}

	if a.VisitId == nil || *a.VisitId <= 0 {
		return errors.New("treatment_visit_id is empty")
	}

	return nil
}

type TreatmentResponse struct {
	Id      uint   `json:"id"`
	Name    string `json:"treatment_name"`
	VisitId uint   `json:"treatment_visit_id"`
}

type TreatmentHistoryResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"treatment_name"`
}
