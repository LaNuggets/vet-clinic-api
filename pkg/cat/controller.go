package cat

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

type CatConfig struct {
	*config.Config
}

func New(configuration *config.Config) *CatConfig {
	return &CatConfig{configuration}
}

// PostHandler godoc
// @Summary      Create a new Cat
// @Description  Creates a new cat entry in the database
// @Tags         cats
// @Accept       json
// @Produce      json
// @Param        cat  body      model.CatRequest  true  "Cat creation payload"
// @Security     BearerAuth
// @Success      200  {object}  model.CatResponse
// @Failure      400  {object}  map[string]string  "Invalid Cat Post request payload"
// @Failure      500  {object}  map[string]string  "Failed to Create specific Cat"
// @Router       /cats [post]
func (config *CatConfig) PostHandler(w http.ResponseWriter, r *http.Request) {

	// Get the request
	req := &model.CatRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid Cat Post request payload. " + err.Error()})
		return
	}

	// Convert the requested data into dbmodel.CatEntry type for the "Create" function
	catEntry := &dbmodel.CatEntry{
		Name:   *req.Name,
		Age:    *req.Age,
		Breed:  *req.Breed,
		Weight: *req.Weight}

	// Request the DB to Create the informations
	entries, err := config.CatEntryRepository.Create(catEntry)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Create specific Cat"})
		return
	}

	// Set up to a dedicated type for the response
	res := &model.CatResponse{
		Id:     entries.ID,
		Name:   entries.Name,
		Age:    entries.Age,
		Breed:  entries.Breed,
		Weight: entries.Weight}

	render.JSON(w, r, res)
}

// GetAllHandler godoc
// @Summary      Get all Cats
// @Description  Find all the cats in the database
// @Tags         cats
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  model.CatResponse
// @Failure      500  {object}  map[string]string  "Failed to retrieve cats"
// @Router       /cats [get]
func (config *CatConfig) GetAllHandler(w http.ResponseWriter, r *http.Request) {

	// Request the DB to get the needed informations
	entries, err := config.CatEntryRepository.FindAll()
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid Find All Cats request payload"})
		return
	}

	// Set up to a dedicated type for the response
	var result []*model.CatResponse
	for _, entrie := range entries {
		result = append(result,
			&model.CatResponse{
				Id:     entrie.ID,
				Name:   entrie.Name,
				Age:    entrie.Age,
				Breed:  entrie.Breed,
				Weight: entrie.Weight})
	}

	render.JSON(w, r, result)
}

// GetByIdHandler godoc
// @Summary      Get cat by ID
// @Description  Retrieves a specific cat from the database by its ID
// @Tags         cats
// @Produce      json
// @Param        id   path      int  true  "Cat ID"
// @Security     BearerAuth
// @Success      200  {object}  model.CatResponse
// @Failure      404  {object}  map[string]string  "Cat not found"
// @Failure      500  {object}  map[string]string  "Failed to find specific cat"
// @Router       /cats/{id} [get]
func (config *CatConfig) GetByIdHandler(w http.ResponseWriter, r *http.Request) {

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

	// Request the DB to get the needed informations
	entries, err := config.CatEntryRepository.FindById(id)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Find specific Cat"})
		return
	}

	// Set up to a dedicated type for the response
	res := &model.CatResponse{
		Id:     entries.ID,
		Name:   entries.Name,
		Age:    entries.Age,
		Breed:  entries.Breed,
		Weight: entries.Weight}

	render.JSON(w, r, res)
}

// GetCatHistoryHandler godoc
// @Summary      Get cat history
// @Description  Retrieves the complete medical history of a cat including visits and treatments
// @Tags         cats
// @Produce      json
// @Param        id   path      int  true  "Cat ID"
// @Security     BearerAuth
// @Success      200  {object}  model.CatHistoryResponse
// @Failure      404  {object}  map[string]string  "Cat not found"
// @Failure      500  {object}  map[string]string  "Failed to find cat history"
// @Router       /cats/{id}/history [get]
func (config *CatConfig) GetCatHistoryHandler(w http.ResponseWriter, r *http.Request) {

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

	// Request the DB to get the needed informations
	entries, err := config.CatEntryRepository.FindCatHistory(id)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Find specific Cat"})
		return
	}

	// Set up to a dedicated type for the response
	var visits []*model.VisitHistoryResponse
	var treatments []*model.TreatmentHistoryResponse

	for _, visit := range entries.Visits {
		for _, treatment := range visit.Treatments {
			treatments = append(treatments, &model.TreatmentHistoryResponse{Id: treatment.ID, Name: treatment.Name})
		}

		visits = append(visits,
			&model.VisitHistoryResponse{
				Id:         visit.ID,
				Date:       visit.Date,
				Reason:     visit.Reason,
				Vet:        visit.Vet,
				Treatments: treatments})
		treatments = nil
	}

	res := &model.CatHistoryResponse{
		Id:     entries.ID,
		Name:   entries.Name,
		Age:    entries.Age,
		Breed:  entries.Breed,
		Weight: entries.Weight,
		Visits: visits}

	render.JSON(w, r, res)
}

// UpdateHandler godoc
// @Summary      Update a cat
// @Description  Updates an existing cat's information in the database
// @Tags         cats
// @Accept       json
// @Produce      json
// @Param        id   path      int                 true  "Cat ID"
// @Param        cat  body      model.CatRequest   true  "Cat update payload"
// @Security     BearerAuth
// @Success      200  {object}  model.CatResponse
// @Failure      400  {object}  map[string]string  "Invalid request payload"
// @Failure      404  {object}  map[string]string  "Cat not found"
// @Failure      500  {object}  map[string]string  "Failed to update cat"
// @Router       /cats/{id} [put]
func (config *CatConfig) UpdateHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the UR
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
	req := &model.CatRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid Cat Update request payload. " + err.Error()})
		return
	}

	// Convert the requested data into dbmodel.CatEntry type for the "Update" function
	catEntry := &dbmodel.CatEntry{
		Name:   *req.Name,
		Age:    *req.Age,
		Breed:  *req.Breed,
		Weight: *req.Weight}

	// Request the DB to Update the informations
	entries, err := config.CatEntryRepository.Update(id, catEntry)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Update Cat"})
		return
	}

	// Set up to a dedicated type for the response
	res := &model.CatResponse{
		Id:     uint(id),
		Name:   entries.Name,
		Age:    entries.Age,
		Breed:  entries.Breed,
		Weight: entries.Weight}

	render.JSON(w, r, res)
}

// DeleteHandler godoc
// @Summary      Delete a cat
// @Description  Deletes a cat from the database by its ID
// @Tags         cats
// @Produce      json
// @Param        id   path      int  true  "Cat ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string  "Cat deleted successfully"
// @Failure      404  {object}  map[string]string  "Cat not found"
// @Failure      500  {object}  map[string]string  "Failed to delete cat"
// @Router       /cats/{id} [delete]
func (config *CatConfig) DeleteHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id in the UR
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
	errDelete := config.CatEntryRepository.DeleteById(id)
	if errDelete != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Delete Cat"})
		return
	}

	render.JSON(w, r, map[string]string{"message": "Cat deleted successfully"})
}
