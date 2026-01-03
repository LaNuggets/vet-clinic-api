package user

import (
	"fmt"
	"net/http"
	"strconv"
	"vet-clinic-api/config"
	"vet-clinic-api/database/dbmodel"
	"vet-clinic-api/pkg/authentication"
	"vet-clinic-api/pkg/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type UserConfig struct {
	*config.Config
}

func New(configuration *config.Config) *UserConfig {
	return &UserConfig{configuration}
}

// LoginHandler godoc
// @Summary      Authenticate a user and get JWT
// @Description  Authenticates a user by email and password, returns a JWT token if credentials are valid.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body model.UserRequest true "Login credentials"
// @Success      200  {object}  model.TokenResponse
// @Failure      400  {object}  map[string]string "Invalid JSON payload"
// @Failure      401  {object}  map[string]string "Invalid email or password"
// @Failure      500  {object}  map[string]string "Failed to generate token"
// @Router       /users/login [post]
func (config *UserConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {

	// Get the request
	req := &model.UserRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid User Post request payload. " + err.Error()})
		return
	}

	// Request the DB to Find the informations
	user, err := config.UserEntryRepository.FindByEmail(*req.Email)
	if err != nil || user == nil {
		render.JSON(w, r, map[string]string{"error": "No user found with email : " + *req.Email})
		return
	}

	// Check User password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*req.Password)) != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate access token for a specific user with 2 hours expiration time
	accessToken, err := authentication.GenerateToken(config.JWTSecret, user.Email, 2)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Generate refresh token for a specific user with 7 days expiration time
	refreshToken, err := authentication.GenerateToken(config.JWTRefreshSecret, user.Email, 7*24)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Set up token to specific response format for better readability
	res := &model.TokensResponse{AccessToken: accessToken, RefreshToken: refreshToken}

	render.JSON(w, r, res)
}

// RefreshHandler godoc
// @Summary      Refresh access token
// @Description  Generate a new access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh body model.RefreshTokenRequest true "Refresh token payload"
// @Success      200  {object}  model.AccessTokenResponse
// @Failure      400  {object}  map[string]string "Invalid JSON payload"
// @Failure      401  {object}  map[string]string "Invalid refresh token"
// @Failure      500  {object}  map[string]string "Failed to generate token"
// @Router       /users/refresh [post]
func (config *UserConfig) RefreshHandler(w http.ResponseWriter, r *http.Request) {

	// Get the request
	req := &model.RefreshTokenRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid refresh token request payload. " + err.Error()})
		return
	}

	// Check the refresh token validity
	email, err := authentication.ParseToken(config.JWTRefreshSecret, *req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Generate new access token with 2 hours expiration time
	newAccessToken, err := authentication.GenerateToken(config.JWTSecret, email, 2)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	res := &model.AccessTokenResponse{AccessToken: newAccessToken}

	render.JSON(w, r, res)
}

// PostHandler godoc
// @Summary      Create a new User
// @Description  Creates a new user entry in the database
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      model.UserRequest  true  "User creation payload"
// @Success      200  {object}  model.UserResponse
// @Failure      400  {object}  map[string]string  "Invalid User Post request payload"
// @Failure      500  {object}  map[string]string  "Failed to Create specific User"
// @Router       /users [post]
func (config *UserConfig) PostHandler(w http.ResponseWriter, r *http.Request) {

	// Get the request
	req := &model.UserRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid User Post request payload. " + err.Error()})
		return
	}

	// Hash the user password for better security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to hash password"})
		return
	}

	// Convert the requested data into dbmodel.UserEntry type for the "Create" function
	userEntry := &dbmodel.UserEntry{Email: *req.Email, Password: string(hashedPassword)}

	// Request the DB to Create the informations
	entries, err := config.UserEntryRepository.Create(userEntry)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Create specific User"})
		return
	}

	// Set up to a dediusered type for the response
	res := &model.UserResponse{
		Id:       entries.ID,
		Email:    entries.Email,
		Password: entries.Password}

	render.JSON(w, r, res)
}

// GetAllHandler godoc
// @Summary      Get all Users
// @Description  Find all the users in the database
// @Tags         users
// @Produce      json
// @Success      200  {array}  model.UserResponse
// @Failure      500  {object}  map[string]string  "Failed to retrieve users"
// @Router       /users [get]
func (config *UserConfig) GetAllHandler(w http.ResponseWriter, r *http.Request) {

	// Request the DB to get the needed informations
	entries, err := config.UserEntryRepository.FindAll()
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid Find All Users request payload"})
		return
	}

	// Set up to a dediusered type for the response
	var result []*model.UserResponse
	for _, entrie := range entries {
		result = append(result,
			&model.UserResponse{
				Id:       entrie.ID,
				Email:    entrie.Email,
				Password: entrie.Password})
	}

	render.JSON(w, r, result)
}

// GetByIdHandler godoc
// @Summary      Get user by ID
// @Description  Retrieves a specific user from the database by its ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  model.UserResponse
// @Failure      404  {object}  map[string]string  "User not found"
// @Failure      500  {object}  map[string]string  "Failed to find specific user"
// @Router       /users/{id} [get]
func (config *UserConfig) GetByIdHandler(w http.ResponseWriter, r *http.Request) {

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
	entries, err := config.UserEntryRepository.FindById(id)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Find specific User"})
		return
	}

	// Set up to a dediusered type for the response
	res := &model.UserResponse{
		Id:       entries.ID,
		Email:    entries.Email,
		Password: entries.Password}

	render.JSON(w, r, res)
}

// UpdateHandler godoc
// @Summary      Update a user
// @Description  Updates an existing user's information in the database
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int                 true  "User ID"
// @Param        user  body      model.UserRequest   true  "User update payload"
// @Success      200  {object}  model.UserResponse
// @Failure      400  {object}  map[string]string  "Invalid request payload"
// @Failure      404  {object}  map[string]string  "User not found"
// @Failure      500  {object}  map[string]string  "Failed to update user"
// @Router       /users/{id} [put]
func (config *UserConfig) UpdateHandler(w http.ResponseWriter, r *http.Request) {

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
	req := &model.UserRequest{}
	if err := render.Bind(r, req); err != nil {
		render.JSON(w, r, map[string]string{"error": "Invalid User Update request payload. " + err.Error()})
		return
	}

	// Convert the requested data into dbmodel.UserEntry type for the "Update" function
	userEntry := &dbmodel.UserEntry{
		Email:    *req.Email,
		Password: *req.Password}

	// Request the DB to Update the informations
	entries, err := config.UserEntryRepository.Update(id, userEntry)
	if err != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Update User"})
		return
	}

	// Set up to a dediusered type for the response
	res := &model.UserResponse{
		Id:       uint(id),
		Email:    entries.Email,
		Password: entries.Password}

	render.JSON(w, r, res)
}

// DeleteHandler godoc
// @Summary      Delete a user
// @Description  Deletes a user from the database by its ID
// @Tags         users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  map[string]string  "User deleted successfully"
// @Failure      404  {object}  map[string]string  "User not found"
// @Failure      500  {object}  map[string]string  "Failed to delete user"
// @Router       /users/{id} [delete]
func (config *UserConfig) DeleteHandler(w http.ResponseWriter, r *http.Request) {

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
	errDelete := config.UserEntryRepository.DeleteById(id)
	if errDelete != nil {
		render.JSON(w, r, map[string]string{"error": "Failed to Delete User"})
		return
	}

	render.JSON(w, r, map[string]string{"message": "User deleted successfully"})
}
