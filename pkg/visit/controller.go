package visit

import (
	"fmt"
	"net/http"
	"strconv"
	"vet-clinic-api/config"
	"vet-clinic-api/database/dbmodel"
	"vet-clinic-api/pkg/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type VisitConfig struct {
	*config.Config
}

func New(configuration *config.Config) *VisitConfig {
	return &VisitConfig{configuration}
}

// PostHandler godoc
// @Summary      Create a new visit
// @Description  Creates a new visit entry in the database
// @Tags         visits
// @Accept       json
// @Produce      json
// @Param        visit  body      model.VisitRequest  true  "Visit creation payload"
// @Security     BearerAuth
// @Success      200    {object}  model.VisitResponse
// @Failure      400    {object}  map[string]string  "Invalid request payload"
// @Failure      500    {object}  map[string]string  "Failed to create visit"
// @Router       /visits [post]
func (config *VisitConfig) PostHandler(w http.ResponseWriter, r *http.Request) {

	// Get the request
	req := &model.VisitRequest{}
	if err := render.Bind(r, req); err != nil {
		// render.JSON(w, r, map[string]string{"error": "Invalid Visit Post request payload"})
		render.JSON(w, r, map[string]string{"error": "Invalid Visit Post request payload. " + err.Error()})
		return
	}

	// Check if the linked cat id existe
	if !config.CatEntryRepository.FindLastCatId(int(*req.CatId)) {
		render.JSON(w, r, map[string]string{"error": "CatId not found in the DB"})
		return
	}

	// Convert the requested data into dbmodel.VisitEntry type for the "Create" function
	visitEntry := &dbmodel.VisitEntry{CatId: *req.CatId, Date: *req.Date, Reason: *req.Reason, Vet: *req.Vet}

	// Request the DB to Create the informations
	entries, err := config.VisitEntryRepository.Create(visitEntry)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Create visit"})
		return
	}

	// Set up to a dedicated type for the response
	res := &model.VisitResponse{
		Id:         entries.ID,
		CatId:      entries.CatId,
		Date:       entries.Date,
		Reason:     entries.Reason,
		Vet:        entries.Vet,
		Treatments: []*model.TreatmentResponse{}}

	render.JSON(w, r, res)
}

// GetAlldHandler godoc
// @Summary      Get all visits with optional filters
// @Description  Retrieves a list of all visits from the database, optionally filtered by vet, reason, or date
// @Tags         visits
// @Produce      json
// @Param        vet     query     string  false  "Filter by veterinarian name"
// @Param        reason  query     string  false  "Filter by visit reason"
// @Param        date    query     string  false  "Filter by date (format: YYYY-MM-DD)"
// @Security     BearerAuth
// @Success      200     {array}   model.VisitHistoryResponse
// @Failure      500     {object}  map[string]string  "Failed to retrieve visits"
// @Router       /visits [get]
func (config *VisitConfig) GetAlldHandler(w http.ResponseWriter, r *http.Request) {

	// Set up the filter by "vet" or "reason"
	var entries []*dbmodel.VisitEntry
	var err error

	vet := r.URL.Query().Get("vet")
	reason := r.URL.Query().Get("reason")
	date := r.URL.Query().Get("date")

	// Request the DB to Get the needed informations base on the filter
	switch {
	case vet != "":
		entries, err = config.VisitEntryRepository.FindByVet(vet)
	case reason != "":
		entries, err = config.VisitEntryRepository.FindByReason(reason)
	case date != "":
		entries, err = config.VisitEntryRepository.FindByDate(date)
	default:
		entries, err = config.VisitEntryRepository.FindAll()
	}

	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Find Visits"})
		return
	}

	// Set up to a dedicated type for the response
	var res []*model.VisitHistoryResponse
	var treatments []*model.TreatmentHistoryResponse

	for _, visit := range entries {
		for _, treatment := range visit.Treatments {
			treatments = append(treatments, &model.TreatmentHistoryResponse{Id: treatment.ID, Name: treatment.Name})
		}

		res = append(res,
			&model.VisitHistoryResponse{
				Id:         visit.ID,
				Date:       visit.Date,
				Reason:     visit.Reason,
				Vet:        visit.Vet,
				Treatments: treatments})
		treatments = nil
	}

	render.JSON(w, r, res)
}

// GetByIdHandler godoc
// @Summary      Get visit by ID
// @Description  Retrieves a specific visit from the database by its ID, including associated treatments
// @Tags         visits
// @Produce      json
// @Param        id   path      int  true  "Visit ID"
// @Security     BearerAuth
// @Success      200  {object}  model.VisitHistoryResponse
// @Failure      404  {object}  map[string]string  "Visit not found"
// @Failure      500  {object}  map[string]string  "Failed to find specific visit"
// @Router       /visits/{id} [get]
func (config *VisitConfig) GetByIdHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error during id convertion")
	}

	// Request the DB to Get the needed informations
	entries, err := config.VisitEntryRepository.FindById(id)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Find specific Visit"})
		return
	}

	// Set up to a dedicated type for the response
	var treatments []*model.TreatmentHistoryResponse

	for _, visit := range entries.Treatments {
		treatments = append(treatments, &model.TreatmentHistoryResponse{Id: visit.ID, Name: visit.Name})
	}

	res := &model.VisitHistoryResponse{
		Id:         entries.ID,
		Date:       entries.Date,
		Reason:     entries.Reason,
		Vet:        entries.Vet,
		Treatments: treatments}

	render.JSON(w, r, res)
}

// UpdateHandler godoc
// @Summary      Update a visit
// @Description  Updates an existing visit's information in the database
// @Tags         visits
// @Accept       json
// @Produce      json
// @Param        id     path      int                  true  "Visit ID"
// @Param        visit  body      model.VisitRequest  true  "Visit update payload"
// @Security     BearerAuth
// @Success      200    {object}  model.VisitResponse
// @Failure      400    {object}  map[string]string  "Invalid request payload"
// @Failure      404    {object}  map[string]string  "Visit not found"
// @Failure      500    {object}  map[string]string  "Failed to update visit"
// @Router       /visits/{id} [put]
func (config *VisitConfig) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error during id convertion")
	}

	// Get the request
	req := &model.VisitRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid Visit Update request payload. " + err.Error()})
		return
	}

	// Check if the linked cat id existe
	if !config.CatEntryRepository.FindLastCatId(int(*req.CatId)) {
		render.JSON(w, r, map[string]string{"error": "CatId not found in the DB"})
		return
	}

	// Convert the requested data into dbmodel.VisitEntry type for the "Update" function
	visitEntry := &dbmodel.VisitEntry{CatId: *req.CatId, Date: *req.Date, Reason: *req.Reason, Vet: *req.Vet}

	// Request the DB to Update the informations
	entries, err := config.VisitEntryRepository.Update(id, visitEntry)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Update Visit"})
		return
	}

	// Set up to a dedicated type for the response
	var treatments []*model.TreatmentResponse
	for _, treatment := range entries.Treatments {
		treatments = append(treatments,
			&model.TreatmentResponse{
				Id:      treatment.ID,
				Name:    treatment.Name,
				VisitId: treatment.VisitId})
	}

	res := &model.VisitResponse{
		Id:         uint(id),
		CatId:      entries.CatId,
		Date:       entries.Date,
		Reason:     entries.Reason,
		Vet:        entries.Vet,
		Treatments: treatments}

	render.JSON(w, r, res)
}

// DeleteHandler godoc
// @Summary      Delete a visit
// @Description  Deletes a visit from the database by its ID
// @Tags         visits
// @Produce      json
// @Param        id   path      int  true  "Visit ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string  "Visit deleted successfully"
// @Failure      404  {object}  map[string]string  "Visit not found"
// @Failure      500  {object}  map[string]string  "Failed to delete visit"
// @Router       /visits/{id} [delete]
func (config *VisitConfig) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the URL
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Error during id convertion")
	}

	// Request the DB to Delete the informations
	errDelete := config.VisitEntryRepository.DeleteById(id)
	if errDelete != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Delete Visit"})
		return
	}

	render.JSON(w, r, map[string]string{"message": "Visit deleted successfully"})
}
