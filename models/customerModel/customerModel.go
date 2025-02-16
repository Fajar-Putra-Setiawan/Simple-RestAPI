package customermodel

import (
	"warehouse-management/config"
	"warehouse-management/entities"
)

func GetAllCustomer() ([]entities.Customer, error) {
	var customers []entities.Customer
	err := config.DB.Find(&customers).Error
	return customers, err
}

func CreateCustomer(customers *entities.Customer) error {
	result := config.DB.Create(customers)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetByID(id uint) (entities.Customer, error) {
	var customer entities.Customer
	err := config.DB.Where("id = ?", id).First(&customer).Error
	if err != nil {
		return customer, err
	}
	return customer, nil
}

func Update(customer *entities.Customer) error {
	return config.DB.Save(customer).Error
}

func Delete(id int) error {
	var customer entities.Customer
	// First, find the category by ID
	if err := config.DB.First(&customer, id).Error; err != nil {
		// If there's another error, return it
		return err
	}

	// Delete the category from the database
	if err := config.DB.Delete(&customer).Error; err != nil {
		// If there's an error during deletion, return the error
		return err
	}

	// If everything goes well, return nil (no error)
	return nil
}
