package entities

import "time"

type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	Address   string    `gorm:"type:text" json:"address"`
	Email     string    `gorm:"type:varchar(100);unique" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
