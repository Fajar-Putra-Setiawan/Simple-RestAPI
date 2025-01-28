package entities

import "time"

// Transaction mencatat informasi terkait transaksi
type Transaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ProductID       uint      `gorm:"not null" json:"product_id"`                        // ID produk yang terlibat dalam transaksi
	Product         Product   `gorm:"foreignKey:ProductID" json:"product"`               // Relasi ke tabel Product
	Quantity        int       `gorm:"not null" json:"quantity"`                          // Jumlah produk yang terjual
	Price           float64   `gorm:"type:decimal(10,2);not null" json:"price"`          // Harga per unit produk
	TotalAmount     float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`   // Total harga transaksi (Quantity * Price)
	TransactionType string    `gorm:"type:varchar(50);not null" json:"transaction_type"` // Jenis transaksi (misal: "sale", "purchase")
	Date            string    `gorm:"type:date;not null" json:"date"`                    // Tanggal transaksi
	CustomerID      uint      `gorm:"not null" json:"customer_id"`                       // ID pelanggan yang melakukan transaksi
	Customer        Customer  `gorm:"foreignKey:CustomerID" json:"customer"`             // Relasi ke tabel Customer
	Status          string    `gorm:"type:varchar(50);default:'pending'" json:"status"`  // Status transaksi (misal: 'completed', 'pending')
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
