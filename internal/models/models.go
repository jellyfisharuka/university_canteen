package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string
	Password string
	Role     string   `gorm:"type:enum('admin', 'user')"`
	Orders   []Order  `gorm:"foreignKey:UserID"`
	Baskets  []Basket `gorm:"foreignKey:UserID"`
}
type Order struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   // Foreign key for User
	OrderStatus  string `gorm:"type:enum('placed', 'preparing', 'ready', 'completed', 'cancelled')"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	TotalPrice   decimal.Decimal
	User         User          `gorm:"foreignKey:UserID"`
	OrderDetails []OrderDetail `gorm:"foreignKey:OrderID"`
}
type OrderDetail struct {
	ID       uint `gorm:"primaryKey;autoIncrement"`
	OrderID  uint // Foreign key for Order
	ItemID   uint // Foreign key for Menu item
	Quantity int
	Order    Order `gorm:"foreignKey:OrderID"`
	MenuItem Menu  `gorm:"foreignKey:ItemID"`
}
type Basket struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint // Foreign key for User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TotalPrice  decimal.Decimal
	User        User         `gorm:"foreignKey:UserID"`
	BasketItems []BasketItem `gorm:"foreignKey:BasketID"`
}
type BasketItem struct {
	ID       uint `gorm:"primaryKey"`
	BasketID uint // Foreign key for Basket
	ItemID   uint // Foreign key for Menu item
	Quantity int
	Price    decimal.Decimal
	Basket   Basket `gorm:"foreignKey:BasketID"`
	MenuItem Menu   `gorm:"foreignKey:ItemID"`
}
type Menu struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Description  string
	Price        decimal.Decimal
	Quantity     int
	IsAvailable  bool
	OrderDetails []OrderDetail `gorm:"foreignKey:ItemID"`
	BasketItems  []BasketItem  `gorm:"foreignKey:ItemID"`
}
