package transactionmodel

import (
	"errors"
	"warehouse-management/config"
	"warehouse-management/entities"
)

func GetAll() ([]entities.Transaction, error) {
	var transaction []entities.Transaction
	result := config.DB.Preload("Product").Find(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func GetProduct() ([]entities.Product, error) {
	var product []entities.Product
	result := config.DB.Find(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

// CreateTransaction menyimpan transaksi dan mengurangi stok produk
func CreateTransaction(transaction *entities.Transaction) error {
	// Mulai transaksi database
	tx := config.DB.Begin()

	// Cek apakah produk ada
	var product entities.Product
	if err := tx.First(&product, transaction.ProductID).Error; err != nil {
		tx.Rollback()
		return errors.New("product not found")
	}

	// Cek apakah stok mencukupi (hanya untuk transaksi "sale")
	if transaction.TransactionType == "sale" && product.Stock < transaction.Quantity {
		tx.Rollback()
		return errors.New("insufficient stock")
	}

	// Kurangi stok produk jika transaksi adalah "sale"
	if transaction.TransactionType == "sale" {
		product.Stock -= transaction.Quantity
	} else if transaction.TransactionType == "purchase" {
		product.Stock += transaction.Quantity
	}

	// Simpan perubahan stok
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to update product stock")
	}

	// Simpan transaksi
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return errors.New("failed to create transaction")
	}

	// Commit transaksi database
	tx.Commit()
	return nil
}

func GetByID(id uint) (entities.Transaction, error) {
	var transaction entities.Transaction
	result := config.DB.Where("ïd = ?", id).Find(&transaction)
	if result.Error != nil {
		return transaction, result.Error
	}
	return transaction, nil
}
