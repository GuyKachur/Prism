package database

import (
	"fmt"
	"refract/refract"

	"github.com/happierall/l"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Datastore Controls interaction with the database & datastore
type Datastore interface {
	StorageAPI
	GetImage(uid string) (*Model, error)
	GetChildren(parentUID string) (*[]Model, error)
	Upload(model *Model) error
	// Update(model *Model) error
	Delete(uid string) error
	LoadImages(page, pageSize int) (*[]Model, error)
	Random() (*Model, error)
	CountTable(table string) (*int64, error)
	Search(term string) (*[]Model, error)

	// //config
	// GetConfig(name string) (*refract.Config, error)
	// CreateConfig(refract.Config) error
}

type instance struct {
	db *gorm.DB
}

//Instance is the datastore client
var Instance Datastore

func init() {
	// dbString := os.Getenv("DB_STRING")
	// fmt.Println(dbString)
	dsn := "user=postgres password=postgres dbname=prism port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Model{}, &refract.Config{}, &refract.Output{})
	if err != nil {
		panic(fmt.Sprintf("Failed to migrate database %v", err))
	}
	// db.Model(&refract.Config{})
	//
	// db.AutoMigrate(&refract.Config{})
	// db.AutoMigrate(&refract.OutputPath{})

	//Set instance
	Instance = &instance{
		db: db,
	}
}

func paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
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
	images := []Model{}
	result := instance.db.Scopes(paginate(page, pageSize)).Find(&images)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "Unable to retrive page: ")
	}
	return &images, nil

}

func (instance *instance) GetImage(uid string) (*Model, error) {
	var model Model
	if uid == "" || uid == "undefined" {
		return nil, errors.New(fmt.Sprintf("Invalid uid: %s", uid))
	}
	if result := instance.db.Where("uid = ?", uid).First(&model); result.Error != nil {
		return nil, errors.Wrap(result.Error, fmt.Sprintf("Unable to retrieve item with UID[%s]: ", uid))
	}
	return &model, nil
}

func (instance *instance) GetChildren(parentUID string) (*[]Model, error) {
	var children = make([]Model, 0)
	if result := instance.db.Where("parent_id != 0 AND parent_id = ?", parentUID).Find(&children); result.Error != nil && result.RowsAffected != 0 {
		return nil, errors.Wrap(result.Error, fmt.Sprintf("Unable to retrieve items with parent_id = %s", parentUID))
	}

	return &children, nil
}
func (instance *instance) Upload(model *Model) error {
	err := model.VerifyUpload()
	if err != nil {
		return err
	}

	err = instance.SaveImage(model)
	if err != nil {
		return errors.Errorf("Error saving image: %v", err)
	}
	if result := instance.db.Create(model); result.Error != nil && result.RowsAffected == 0 {
		return result.Error
	}
	l.Debug(model.UID)

	return nil
}

// Not sure we need an update yet, left updated time in model, but,,, ehh
// func (instance *instance) Update(model *Model) error {
// 	return nil
// }
func (instance *instance) Delete(uid string) error {
	if result := instance.db.Delete(&Model{}, uid); result.Error != nil {
		return errors.Wrap(result.Error, fmt.Sprintf("Error deleting model: %s", uid))
	}
	return nil
}

func (instance *instance) CountTable(table string) (*int64, error) {
	var count int64
	if result := instance.db.Table(table).Count(&count); result.Error != nil {
		return nil, errors.Wrap(result.Error, fmt.Sprintf("Table %s unaccesible", table))
	}
	return &count, nil
}
func (instance *instance) Random() (*Model, error) {
	model := Model{}
	if result := instance.db.Limit(1).Order("RANDOM()").Find(&model); result.Error != nil && result.RowsAffected != 0 {
		return nil, result.Error
	}
	return &model, nil

}

func (instance *instance) Search(term string) (*[]Model, error) {
	var images = make([]Model, 0)
	if err := instance.db.Where("tags LIKE %?%", term).Find(&images).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Unable to retrieve items for term: %s", term))
	}
	return &images, nil

}

// func (instance *instance) GetConfig(name string) (*refract.Config, error) {
// 	return nil, nil
// }
