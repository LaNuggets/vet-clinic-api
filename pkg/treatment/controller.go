package treatment

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

type TreatmentConfig struct {
	*config.Config
}

func New(configuration *config.Config) *TreatmentConfig {
	return &TreatmentConfig{configuration}
}

// PostHandler godoc
// @Summary      Create a new treatment
// @Description  Creates a new treatment entry in the database
// @Tags         treatments
// @Accept       json
// @Produce      json
// @Param        treatment  body      model.TreatmentRequest  true  "Treatment creation payload"
// @Success      200        {object}  model.TreatmentResponse
// @Failure      400        {object}  map[string]string  "Invalid request payload"
// @Failure      500        {object}  map[string]string  "Failed to create treatment"
// @Router       /treatments [post]
func (config *TreatmentConfig) PostHandler(w http.ResponseWriter, r *http.Request) {

	// Get the request
	req := &model.TreatmentRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid Treatment Post request payload. " + err.Error()})
		return
	}

	// Check if the linked cat id existe
	if !config.VisitEntryRepository.FindLastVisitId(int(*req.VisitId)) {
		render.JSON(w, r, map[string]string{"error": "VisitId not found in the DB"})
		return
	}

	// Convert the requested data into dbmodel.TreatmentEntry type for the "Create" function
	treatmentEntry := &dbmodel.TreatmentEntry{Name: *req.Name, VisitId: uint(*req.VisitId)}

	// Request the DB to Create the informations
	entries, err := config.TreatmentEntryRepository.Create(treatmentEntry)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Create Treatment"})
		return
	}

	// Set up to a dedicated type for the response
	res := &model.TreatmentResponse{Id: entries.ID, Name: entries.Name, VisitId: entries.VisitId}

	render.JSON(w, r, res)
}

// GetAllHandler godoc
// @Summary      Get all treatments
// @Description  Retrieves a list of all treatments from the database
// @Tags         treatments
// @Produce      json
// @Success      200  {array}   model.TreatmentResponse
// @Failure      500  {object}  map[string]string  "Failed to retrieve treatments"
// @Router       /treatments [get]
func (config *TreatmentConfig) GetAllHandler(w http.ResponseWriter, r *http.Request) {

	// Request the DB to Get the needed informations
	entries, err := config.TreatmentEntryRepository.FindAll()
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Find All treatments"})
		return
	}

	// Set up to a dedicated type for the response
	var result []*model.TreatmentResponse
	for _, entrie := range entries {
		result = append(result, &model.TreatmentResponse{Id: entrie.ID, Name: entrie.Name})
	}

	render.JSON(w, r, result)
}

// GetByIdHandler godoc
// @Summary      Get treatment by ID
// @Description  Retrieves a specific treatment from the database by its ID
// @Tags         treatments
// @Produce      json
// @Param        id   path      int  true  "Treatment ID"
// @Success      200  {object}  model.TreatmentResponse
// @Failure      404  {object}  map[string]string  "Treatment not found"
// @Failure      500  {object}  map[string]string  "Failed to find specific treatment"
// @Router       /treatments/{id} [get]
func (config *TreatmentConfig) GetByIdHandler(w http.ResponseWriter, r *http.Request) {

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
	entries, err := config.TreatmentEntryRepository.FindById(id)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Find specific Treatment"})
		return
	}

	// Set up to dedicated type for the response
	res := &model.TreatmentResponse{Id: entries.ID, Name: entries.Name, VisitId: entries.VisitId}

	render.JSON(w, r, res)
}

// GetByVisitIdHandler godoc
// @Summary      Get treatments by visit ID
// @Description  Retrieves all treatments associated with a specific visit
// @Tags         treatments
// @Produce      json
// @Param        id   path      int  true  "Visit ID"
// @Success      200  {array}   model.TreatmentHistoryResponse
// @Failure      404  {object}  map[string]string  "Visit not found"
// @Failure      500  {object}  map[string]string  "Failed to find treatments for visit"
// @Router       /treatments/{id}/history [get]
func (config *TreatmentConfig) GetByVisitIdHandler(w http.ResponseWriter, r *http.Request) {

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
	entries, err := config.TreatmentEntryRepository.FindByVisitId(id)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Find All treatments for a specific visit"})
		return
	}

	// Set up to a dedicated type for the response
	var result []*model.TreatmentHistoryResponse
	for _, entrie := range entries {
		res := &model.TreatmentHistoryResponse{Id: entrie.ID, Name: entrie.Name}
		result = append(result, res)
	}

	render.JSON(w, r, result)
}

// UpdateHandler godoc
// @Summary      Update a treatment
// @Description  Updates an existing treatment's information in the database
// @Tags         treatments
// @Accept       json
// @Produce      json
// @Param        id         path      int                      true  "Treatment ID"
// @Param        treatment  body      model.TreatmentRequest  true  "Treatment update payload"
// @Success      200        {object}  model.TreatmentResponse
// @Failure      400        {object}  map[string]string  "Invalid request payload"
// @Failure      404        {object}  map[string]string  "Treatment not found"
// @Failure      500        {object}  map[string]string  "Failed to update treatment"
// @Router       /treatments/{id} [put]
func (config *TreatmentConfig) UpdateHandler(w http.ResponseWriter, r *http.Request) {

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
	req := &model.TreatmentRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid Treatment Update request payload. " + err.Error()})
		return
	}

	// Check if the linked cat id existe
	if !config.VisitEntryRepository.FindLastVisitId(int(*req.VisitId)) {
		render.JSON(w, r, map[string]string{"error": "VisitId not found in the DB"})
		return
	}

	// Convert the requested data into dbmodel.TreatmentEntry type for the "Update" function
	treatmentEntry := &dbmodel.TreatmentEntry{Name: *req.Name, VisitId: *req.VisitId}

	// Request the DB to Update the informations
	entries, err := config.TreatmentEntryRepository.Update(id, treatmentEntry)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Update Treatment"})
		return
	}

	// Set up to a dedicated type for the response
	res := &model.TreatmentResponse{Id: uint(id), Name: entries.Name, VisitId: entries.VisitId}

	render.JSON(w, r, res)
}

// DeleteHandler godoc
// @Summary      Delete a treatment
// @Description  Deletes a treatment from the database by its ID
// @Tags         treatments
// @Produce      json
// @Param        id   path      int  true  "Treatment ID"
// @Success      200  {object}  map[string]string  "Treatment deleted successfully"
// @Failure      404  {object}  map[string]string  "Treatment not found"
// @Failure      500  {object}  map[string]string  "Failed to delete treatment"
// @Router       /treatments/{id} [delete]
func (config *TreatmentConfig) DeleteHandler(w http.ResponseWriter, r *http.Request) {

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
	errDelete := config.TreatmentEntryRepository.DeleteById(id)
	if errDelete != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Delete Treatment"})
		return
	}

	render.JSON(w, r, map[string]string{"message": "Treatment deleted successfully"})
}
