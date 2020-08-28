package database

import (
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Datastore interface {
	GetImage(uid string) (*Model, error)
	GetChildren(parentUID string) (*[]Model, error)
	Upload(model *Model) error
	// Update(model *Model) error
	Delete(uid string) error
	LoadImages(page, pageSize int) (*[]Model, error)
}

type instance struct {
	db *gorm.DB
}

var Instance Datastore

func init() {
	// dbString := os.Getenv("DB_STRING")
	// fmt.Println(dbString)
	dsn := "user=refract password=postgres dbname=refract port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&Model{})

	//Set instance
	Instance = &instance{
		db: db,
	}
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (instance *instance) LoadImages(page, pageSize int) (*[]Model, error) {
	var images = make([]Model, 0)
	if err := instance.db.Scopes(Paginate(page, pageSize)).Find(&Model{}).Error; err != nil {
		return nil, errors.Wrap(err, "Unable to retrive page: ")
	}
	return &images, nil

}

func (instance *instance) GetImage(uid string) (*Model, error) {
	var model Model
	if result := instance.db.Where("UID = ?", uid).First(&model); result.Error != nil {
		return nil, errors.Wrap(result.Error, fmt.Sprintf("Unable to retrieve item with UID: %s", uid))
	}
	return &model, nil
}

func (instance *instance) GetChildren(parentUID string) (*[]Model, error) {
	var children = make([]Model, 0)
	if err := instance.db.Where("parentUID <> ?", parentUID).Find(&children).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unable to retrieve items with parent: %s", parentUID))
	}
	return &children, nil
}
func (instance *instance) Upload(model *Model) error {
	err := model.Verify()
	if err != nil {
		return err
	}
	if result := instance.db.Create(model); result.Error != nil {
		return errors.Wrap(err, "Error creating model in database: ")
	}
	return nil
}

// Not sure we need an update yet, left updated time in model, but,,, ehh
// func (instance *instance) Update(model *Model) error {
// 	return nil
// }
func (instance *instance) Delete(uid string) error {
	if result := instance.db.Delete(&Model{}, uid); result.Error != nil {
		return errors.Wrap(result.Error, fmt.Sprintf("Error deleting model: %d", uid))
	}
	return nil
}
