package models

import (
	"bda/utils"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func FormatDbDsn() string {
	host := utils.Getenv("DB_HOST", "172.22.0.2")
	port := utils.Getenv("DB_PORT", "5432")
	dbName := utils.Getenv("DB_NAME", "nodes")
	acc := utils.Getenv("DB_USER", "nodes_manager")
	pwd := utils.Getenv("DB_USER_PWD", "nodes_manager_pwd")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, acc, pwd, dbName, port)
}

func GetDb() (*gorm.DB, error) {
	config := gorm.Config{}
	db, err := gorm.Open(postgres.Open(FormatDbDsn()), &config)
	if err != nil {
		return &gorm.DB{}, err
	}

	// Migrate the DB if possible
	db.AutoMigrate(Node{})
	return db, nil
}
