package dbmodel

import (
	"gorm.io/gorm"
)

type UserEntry struct {
	gorm.Model
	Email    string `json:"user_email"`
	Password string `json:"user_password"`
	Role     string `json:"user_role"`
}

type UserEntryRepository interface {
	Create(entry *UserEntry) (*UserEntry, error)
	FindAll() ([]*UserEntry, error)
	FindById(id int) (*UserEntry, error)
	FindByEmail(email string) (*UserEntry, error)
	Update(id int, entry *UserEntry) (*UserEntry, error)
	DeleteById(id int) error
}

type userEntryRepository struct {
	db *gorm.DB
}

func NewUserEntryRepository(db *gorm.DB) UserEntryRepository {
	return &userEntryRepository{db: db}
}

func (r *userEntryRepository) Create(entry *UserEntry) (*UserEntry, error) {

	if err := r.db.Create(entry).Error; err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *userEntryRepository) FindAll() ([]*UserEntry, error) {

	var entries []*UserEntry
	if err := r.db.Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *userEntryRepository) FindById(id int) (*UserEntry, error) {

	var entries *UserEntry
	if err := r.db.Model(&UserEntry{}).First(&entries, id).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *userEntryRepository) FindByEmail(email string) (*UserEntry, error) {

	var entries *UserEntry
	if err := r.db.Where("email = ?", email).First(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *userEntryRepository) Update(id int, entry *UserEntry) (*UserEntry, error) {

	result := r.db.Model(&UserEntry{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"email":    entry.Email,
			"password": entry.Password,
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

func (r *userEntryRepository) DeleteById(id int) error {

	if err := r.db.Delete(&UserEntry{}, id).Error; err != nil {
		return err
	}

	return nil
}
