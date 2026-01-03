package config

import (
	"os"
	"vet-clinic-api/database"
	"vet-clinic-api/database/dbmodel"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {

	// Access token
	JWTSecret string

	// Refresh token
	JWTRefreshSecret string

	// Repository connection
	CatEntryRepository       dbmodel.CatEntryRepository
	TreatmentEntryRepository dbmodel.TreatmentEntryRepository
	VisitEntryRepository     dbmodel.VisitEntryRepository
	UserEntryRepository      dbmodel.UserEntryRepository
}

func New() (*Config, error) {

	config := Config{}

	// Init DB connection
	databaseSession, err := gorm.Open(sqlite.Open("vet_clinic_api.db"), &gorm.Config{})
	if err != nil {
		return &config, err
	}

	config.JWTSecret = os.Getenv("JWT_SECRET")
	config.JWTRefreshSecret = os.Getenv("JWT_REFRESH_SECRET")

	// Models migrate
	database.Migrate(databaseSession)

	// Init repository
	config.CatEntryRepository = dbmodel.NewCatEntryRepository(databaseSession)
	config.TreatmentEntryRepository = dbmodel.NewTreatmentEntryRepository(databaseSession)
	config.VisitEntryRepository = dbmodel.NewVisitEntryRepository(databaseSession)
	config.UserEntryRepository = dbmodel.NewUserEntryRepository(databaseSession)

	return &config, nil
}
