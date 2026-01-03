package dbmodel

import (
	"gorm.io/gorm"
)

type TreatmentEntry struct {
	gorm.Model
	Name    string `json:"treatment_name"`
	VisitId uint   `json:"treatment_visit_id"`
}

type TreatmentEntryRepository interface {
	Create(entry *TreatmentEntry) (*TreatmentEntry, error)
	FindAll() ([]*TreatmentEntry, error)
	FindByVisitId(id int) ([]*TreatmentEntry, error)
	FindById(id int) (*TreatmentEntry, error)
	Update(id int, entry *TreatmentEntry) (*TreatmentEntry, error)
	DeleteById(id int) error
}

type treatmentEntryRepository struct {
	db *gorm.DB
}

func NewTreatmentEntryRepository(db *gorm.DB) TreatmentEntryRepository {
	return &treatmentEntryRepository{db: db}
}

func (r *treatmentEntryRepository) Create(entry *TreatmentEntry) (*TreatmentEntry, error) {

	if err := r.db.Create(entry).Error; err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *treatmentEntryRepository) FindAll() ([]*TreatmentEntry, error) {

	var entries []*TreatmentEntry
	if err := r.db.Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *treatmentEntryRepository) FindByVisitId(id int) ([]*TreatmentEntry, error) {

	var entries []*TreatmentEntry
	if err := r.db.Where("visit_id = ?", id).Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *treatmentEntryRepository) FindById(id int) (*TreatmentEntry, error) {

	var entries *TreatmentEntry
	if err := r.db.First(&entries, id).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *treatmentEntryRepository) Update(id int, entry *TreatmentEntry) (*TreatmentEntry, error) {

	result := r.db.Model(&TreatmentEntry{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":     entry.Name,
			"visit_id": entry.VisitId,
		})

	if result.Error != nil {
		return nil, result.Error
	}

	// Check if something has been update
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return entry, nil
}

func (r *treatmentEntryRepository) DeleteById(id int) error {

	if err := r.db.Delete(&TreatmentEntry{}, id).Error; err != nil {
		return err
	}

	return nil
}
