package categoriesmodel

import (
	"warehouse-management/config"
	"warehouse-management/entities"
)

func GetAll() ([]entities.Category, error) {
	var category []entities.Category
	err := config.DB.Find(&category).Error
	return category, err
}

func Create(category *entities.Category) error {
	// Gunakan metode Create untuk menyimpan data ke database
	result := config.DB.Create(category)

	// Jika ada error saat menyimpan, kembalikan error-nya
	if result.Error != nil {
		return result.Error
	}

	// Jika berhasil, kembalikan nil (tidak ada error)
	return nil
}

func Update(category *entities.Category) error {
	return config.DB.Save(category).Error
}

func GetByID(id uint) (entities.Category, error) {
	var category entities.Category
	err := config.DB.Where("id = ?", id).First(&category).Error
	if err != nil {
		return category, err
	}
	return category, nil
}
