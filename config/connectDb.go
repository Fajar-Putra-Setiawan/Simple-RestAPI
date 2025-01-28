package config

import (
	"log"
	"warehouse-management/entities"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/go_wms_restapi?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Successfully connected to the database")

	err = DB.AutoMigrate(&entities.Product{}, &entities.Category{}, &entities.Customer{}, &entities.Transaction{})
	if err != nil {
		panic("Failed to migrate database" + err.Error())
	}

	log.Println("Database migration completed successfully")
}
