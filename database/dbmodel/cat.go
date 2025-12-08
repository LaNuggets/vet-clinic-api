package dbmodel

import (
	"gorm.io/gorm"
)

type CatEntry struct {
	gorm.Model
	Name   string `json:"cat_name"`
	Age    int    `json:"cat_age"`
	Breed  string `json:"cat_breed"`
	Weight int    `json:"cat_weight"`

	//Add a foreignKey to CatId on the table Visit, and a Delete On Cascade
	Visits []VisitEntry `json:"visits" gorm:"foreignKey:CatId; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;"`
}

type CatEntryRepository interface {
	Create(entry *CatEntry) (*CatEntry, error)
	FindAll() ([]*CatEntry, error)
	FindById(id int) (*CatEntry, error)
	FindCatHistory(id int) (*CatEntry, error)
	FindLastCatId(id int) bool
	Update(id int, entry *CatEntry) (*CatEntry, error)
	DeleteById(id int) error
}

type catEntryRepository struct {
	db *gorm.DB
}

func NewCatEntryRepository(db *gorm.DB) CatEntryRepository {
	return &catEntryRepository{db: db}
}

func (r *catEntryRepository) Create(entry *CatEntry) (*CatEntry, error) {

	if err := r.db.Create(entry).Error; err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *catEntryRepository) FindAll() ([]*CatEntry, error) {

	var entries []*CatEntry
	if err := r.db.Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *catEntryRepository) FindById(id int) (*CatEntry, error) {

	var entries *CatEntry
	if err := r.db.Model(&CatEntry{}).
		Preload("Visits.Treatments").
		First(&entries, id).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *catEntryRepository) FindCatHistory(id int) (*CatEntry, error) {

	var entries *CatEntry
	if err := r.db.Model(&CatEntry{}).
		Preload("Visits.Treatments").
		First(&entries, id).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *catEntryRepository) FindLastCatId(id int) bool {

	var count int64
	r.db.Model(&CatEntry{}).Where("id = ?", id).Count(&count)

	return count > 0
}

func (r *catEntryRepository) Update(id int, entry *CatEntry) (*CatEntry, error) {

	if err := r.db.Model(&CatEntry{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":   entry.Name,
			"age":    entry.Age,
			"breed":  entry.Breed,
			"weight": entry.Weight,
		}).Error; err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *catEntryRepository) DeleteById(id int) error {

	if err := r.db.Delete(&CatEntry{}, id).Error; err != nil {
		return err
	}

	return nil
}
