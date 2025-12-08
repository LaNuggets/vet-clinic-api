package config

import (
	"vet-clinic-api/database"
	"vet-clinic-api/database/dbmodel"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	// Repository connection
	CatEntryRepository       dbmodel.CatEntryRepository
	TreatmentEntryRepository dbmodel.TreatmentEntryRepository
	VisitEntryRepository     dbmodel.VisitEntryRepository
}

func New() (*Config, error) {

	config := Config{}

	// Init DB connection
	databaseSession, err := gorm.Open(sqlite.Open("vet_clinic_api.db"), &gorm.Config{})
	if err != nil {
		return &config, err
	}

	// Models migrate
	database.Migrate(databaseSession)

	// Init repository
	config.CatEntryRepository = dbmodel.NewCatEntryRepository(databaseSession)
	config.TreatmentEntryRepository = dbmodel.NewTreatmentEntryRepository(databaseSession)
	config.VisitEntryRepository = dbmodel.NewVisitEntryRepository(databaseSession)

	return &config, nil
}
