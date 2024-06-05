package db

import (
	"log"
	"os"
	"todo-go/internal/store/dbstore"

	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

func opensqlite(dbName string) (*gorm.DB, error) {

	// make the temp directory if it doesn't exist
	err := os.MkdirAll("/tmp", 0755)
	if err != nil {
		return nil, err
	}

	return gorm.Open(sqlite.Open(dbName+".db"), &gorm.Config{})
}

func MustOpenSqlite(dbName string) *gorm.DB {

	db, err := opensqlite(dbName)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&dbstore.User{}, &dbstore.Session{}, &dbstore.Todo{}, &dbstore.Role{})

	if err != nil {
		panic(err)
	}
	log.Default().Print("sqlite db created")
	return db
}
