package db

import (
	"github.com/glebarez/sqlite"
	"github.com/samber/do/v2"
	"gorm.io/gorm"
)

func NewDatabase(i do.Injector) (*Database, error) {
    return &Database{}, nil
}

type Database struct {
    Database *gorm.DB
}

func (database *Database) Initialize(path string) (error) {
    db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})

    if err != nil {
        return err
    }

    database.Database = db

    return nil
}

func (database *Database) Migrate() error {
    return database.Database.AutoMigrate(&UsersTable{}, &LanguagesTable{})
}
