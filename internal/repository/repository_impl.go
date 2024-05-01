package repository

import (
	"final_project/internal/models"

	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(user models.User) error {
	return ur.db.Create(&user).Error
}

func (ur *userRepository) GetUserByID(id uint) (models.User, error) {
	var user models.User
	result := ur.db.First(&user, id)
	return user, result.Error
}

func (ur *userRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	result := ur.db.Where("username = ?", username).First(&user)
	return user, result.Error
}
func NewOrderRepository(db *gorm.DB) OrderRepository {
    return &orderRepository{db: db}
}

func (or *orderRepository) CreateOrder(order models.Order) error {
    return or.db.Create(&order).Error
}

func (or *orderRepository) GetOrderByID(id uint) (models.Order, error) {
    var order models.Order
    result := or.db.First(&order, id)
    return order, result.Error
}
func NewMenuRepository(db *gorm.DB) MenuRepository {
    return &menuRepository{db: db}
}

func (mr *menuRepository) CreateMenuItem(menuItem models.Menu) error {
    return mr.db.Create(&menuItem).Error
}

func (mr *menuRepository) GetMenuItemByID(id uint) (models.Menu, error) {
    var menuItem models.Menu
    result := mr.db.First(&menuItem, id)
    return menuItem, result.Error
}

func (mr *menuRepository) UpdateMenuItem(menuItem models.Menu) error {
    return mr.db.Save(&menuItem).Error
}

func (mr *menuRepository) DeleteMenuItem(id uint) error {
    return mr.db.Delete(&models.Menu{}, id).Error
}

func (mr *menuRepository) GetAllMenuItems() ([]models.Menu, error) {
    var menuItems []models.Menu
    result := mr.db.Find(&menuItems)
    return menuItems, result.Error
}