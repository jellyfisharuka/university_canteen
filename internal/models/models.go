package models

import (
	"errors"
	"gorm.io/gorm"
	"time"

	"github.com/shopspring/decimal"
)

type Role string

const (
	Admin  Role = "admin"
	Client Role = "client"
)

type Status string

const (
	Canceled   Status = "canseled"
	Preparing  Status = "preparing"
	Ready      Status = "ready"
	Compleated Status = "compoleated"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string
	Password string
	Role     Role
	Orders   []Order  `gorm:"foreignKey:UserID"`
	Baskets  []Basket `gorm:"foreignKey:UserID"`
}
type Order struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   // Foreign key for User
	OrderStatus  Status `gorm:"type:varchar(255)"`
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
	TotalCost    decimal.Decimal
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

func (o *Order) BeforeSave(tx *gorm.DB) (err error) {
	switch o.OrderStatus {
	case Canceled, Preparing, Ready, Compleated:
		return nil
	default:
		return errors.New("invalid order status")
	}
}
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	switch u.Role {
	case Admin, Client:
		return nil
	default:
		return errors.New("invalid user role")
	}
}
