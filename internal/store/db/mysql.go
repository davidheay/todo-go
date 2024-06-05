package db

import (
	"fmt"
	"log"
	"todo-go/internal/store/dbstore"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func openMysql(dbName string, password string) (*gorm.DB, error) {
	dataSourceName := fmt.Sprintf("root:%s@tcp(localhost:3306)/%s", password, dbName)
	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil

}

func MustOpenMysql(dbName string, password string) (*gorm.DB, error) {

	db, err := openMysql(dbName, password)
	if err != nil {
		return nil, err
	}
	// db.Migrator().DropTable(&dbstore.User{}, &dbstore.Session{}, &dbstore.Todo{}, &dbstore.Role{})
	err = db.AutoMigrate(&dbstore.User{}, &dbstore.Session{}, &dbstore.Todo{}, &dbstore.Role{})

	if err != nil {
		return nil, err
	}
	log.Default().Println("mysql db created")

	return db, nil
}
