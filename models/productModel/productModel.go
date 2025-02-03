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

func Create(product *entities.Product) error {
	// Gunakan metode Create untuk menyimpan data ke database
	result := config.DB.Create(product)

	// Jika ada error saat menyimpan, kembalikan error-nya
	if result.Error != nil {
		return result.Error
	}

	// Jika berhasil, kembalikan nil (tidak ada error)
	return nil
}

func Update(product *entities.Product) error {
	return config.DB.Save(product).Error
}

// GetCategories retrieves all categories from the database
func GetCategories() ([]entities.Category, error) {
	var categories []entities.Category
	result := config.DB.Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

// CategoryExists checks if a category exists in the database
func CategoryExists(categoryID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&entities.Category{}).Where("id = ?", categoryID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetByID(id uint) (entities.Product, error) {
	var product entities.Product
	err := config.DB.Where("id = ?", id).First(&product).Error
	if err != nil {
		return product, err
	}
	return product, nil
}

func Delete(id int) error {
	var product entities.Product
	// First, find the category by ID
	if err := config.DB.First(&product, id).Error; err != nil {
		// If there's another error, return it
		return err
	}

	// Delete the category from the database
	if err := config.DB.Delete(&product).Error; err != nil {
		// If there's an error during deletion, return the error
		return err
	}

	// If everything goes well, return nil (no error)
	return nil
}
