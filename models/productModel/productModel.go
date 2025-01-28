package productmodel

import (
	"warehouse-management/config"
	"warehouse-management/entities"
)

func GetAll() ([]entities.Product, error) {
	var products []entities.Product
	err := config.DB.Preload("Category").Find(&products).Error
	return products, err
}
