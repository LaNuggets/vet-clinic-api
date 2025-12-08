package dbmodel

import (
	"gorm.io/gorm"
)

type VisitEntry struct {
	gorm.Model
	CatId  uint   `json:"visit_cat_id"`
	Date   string `json:"visit_date"`
	Reason string `json:"visit_reason"`
	Vet    string `json:"visit_vet"`

	//Add a foreignKey to VisitId on the table Treatment, and a Delete On Cascade
	Treatments []TreatmentEntry `json:"treatments" gorm:"foreignKey:VisitId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type VisitEntryRepository interface {
	Create(entry *VisitEntry) (*VisitEntry, error)
	FindAll() ([]*VisitEntry, error)
	FindById(id int) (*VisitEntry, error)
	FindByReason(reason string) ([]*VisitEntry, error)
	FindByVet(vet string) ([]*VisitEntry, error)
	FindByDate(date string) ([]*VisitEntry, error)
	FindLastVisitId(id int) bool
	Update(id int, entry *VisitEntry) (*VisitEntry, error)
	DeleteById(id int) error
}

type visitEntryRepository struct {
	db *gorm.DB
}

func NewVisitEntryRepository(db *gorm.DB) VisitEntryRepository {
	return &visitEntryRepository{db: db}
}

func (r *visitEntryRepository) Create(entry *VisitEntry) (*VisitEntry, error) {

	if err := r.db.Create(entry).Error; err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *visitEntryRepository) FindAll() ([]*VisitEntry, error) {

	var entries []*VisitEntry
	if err := r.db.Model(&VisitEntry{}).
		Preload("Treatments").
		Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *visitEntryRepository) FindById(id int) (*VisitEntry, error) {

	var entries *VisitEntry
	if err := r.db.Model(&VisitEntry{}).
		Preload("Treatments").
		First(&entries, id).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *visitEntryRepository) FindByReason(reason string) ([]*VisitEntry, error) {

	var entries []*VisitEntry
	if err := r.db.Model(&VisitEntry{}).
		Preload("Treatments").
		Where("reason LIKE ?", "%"+reason+"%").
		Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *visitEntryRepository) FindByVet(vet string) ([]*VisitEntry, error) {

	var entries []*VisitEntry
	if err := r.db.Model(&VisitEntry{}).
		Preload("Treatments").
		Where("vet LIKE ?", "%"+vet+"%").
		Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *visitEntryRepository) FindByDate(date string) ([]*VisitEntry, error) {

	var entries []*VisitEntry
	if err := r.db.Model(&VisitEntry{}).
		Preload("Treatments").
		Where("DATE(date) = ?", date).
		Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *visitEntryRepository) FindLastVisitId(id int) bool {

	var count int64
	r.db.Model(&VisitEntry{}).Where("id = ?", id).Count(&count)

	return count > 0
}

func (r *visitEntryRepository) Update(id int, entry *VisitEntry) (*VisitEntry, error) {

	if err := r.db.Model(&VisitEntry{}).
		Preload("Treatments").
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"cat_id": entry.CatId,
			"date":   entry.Date,
			"reason": entry.Reason,
			"vet":    entry.Vet,
		}).Error; err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *visitEntryRepository) DeleteById(id int) error {

	if err := r.db.Delete(&VisitEntry{}, id).Error; err != nil {
		return err
	}

	return nil
}
