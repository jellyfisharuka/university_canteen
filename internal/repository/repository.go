package repository

import (
	"final_project/internal/models"

	"gorm.io/gorm"
)

type MenuRepository interface {
	CreateMenuItem(menuItem models.Menu) error
	GetMenuItemByID(id uint) (models.Menu, error)
	UpdateMenuItem(menuItem models.Menu) error
	DeleteMenuItem(id uint) error
	GetAllMenuItems() ([]models.Menu, error)
}
type menuRepository struct {
    db *gorm.DB
}

type UserRepository interface {
    CreateUser(user models.User) error
    GetUserByID(id uint) (models.User, error)
    GetUserByUsername(username string) (models.User, error)
}
type userRepository struct {
    db *gorm.DB
}
type OrderRepository interface {
    CreateOrder(order models.Order) error
    GetOrderByID(id uint) (models.Order, error)
}
type orderRepository struct {
    db *gorm.DB
}



