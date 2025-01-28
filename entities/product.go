package entities

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(300)" json:"nama_product"`
	Description string    `gorm:"type:text" json:"deskripsi"`
	Stock       int       `gorm:"type:int" json:"stok"`
	CategoryID  uint      `gorm:"not null" json:"category_id"`           // Relasi ke tabel Category
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category"` // Relasi ke tabel Category
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
