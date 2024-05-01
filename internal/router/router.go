package router

import (
	"final_project/initializers"
	"final_project/internal/models"
	"final_project/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	setupPublicEndpoints(router)
	setupPrivateEndpoints(router)
	setupAuthEndpoints(router)
	setupBasketEndpoints(router)
	setupMenuEndpoints(router)
	setupOrderEndpoints(router)
	return router
}

func setupPublicEndpoints(router *gin.Engine) {
	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "welcome to public endpoint"})
	})
}

func setupPrivateEndpoints(router *gin.Engine) {
	router.GET("/private", utils.AuthMiddleware(), func(c *gin.Context) {
		role, _ := c.Get("role") // Получаем роль пользователя из контекста
		if role == "admin" {
			// Действия для администратора меню
			c.JSON(http.StatusOK, gin.H{"message": "welcome to private endpoint (menu admin)"})
		} else if role == "client" {
			// Действия для обычного пользователя
			c.JSON(http.StatusOK, gin.H{"message": "welcome to private endpoint (user)"})
		}

	})
}

func setupAuthEndpoints(router *gin.Engine) {
	router.POST("/login", func(c *gin.Context) {
		var loginUser models.User
		if err := c.BindJSON(&loginUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		var existingUser models.User
		result := initializers.DB.Select("ID", "username", "password", "role").Where("username = ?", loginUser.Username).First(&existingUser)
		if result.Error != nil || !utils.CheckPassword(existingUser.Password, loginUser.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		token, err := utils.GenerateToken(existingUser.Username, string(existingUser.Role), int(existingUser.ID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		fmt.Println("Generated token:", token)
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
	router.POST("/signup", func(c *gin.Context) {
		var newUser models.User
		if err := c.BindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := utils.SignupUser(initializers.DB, newUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign up user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User signed up successfully"})
	})

}

func setupMenuEndpoints(router *gin.Engine) {
	menuRoutes := router.Group("/menu", utils.AuthMiddleware())
	{
		menuRoutes.POST("/", func(c *gin.Context) {
			role, _ := c.Get("role")
			if role == "admin" {
				var menuItem models.Menu
				if err := c.BindJSON(&menuItem); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
					return
				}
				if err := initializers.DB.Create(&menuItem).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add menu item"})
					return
				}
				c.JSON(http.StatusCreated, gin.H{"message": "Menu item added successfully", "menuItemId": menuItem.ID})
			} else if role == "client" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
				return
			}
		})

		menuRoutes.PUT("/:itemId", func(c *gin.Context) {
			role, _ := c.Get("role")
			if role != "admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
				return
			}
			var menuItem models.Menu
			itemId := c.Param("itemId")
			if err := initializers.DB.First(&menuItem, itemId).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Menu item not found"})
				return
			}
			if err := c.BindJSON(&menuItem); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}
			initializers.DB.Save(&menuItem)
			c.JSON(http.StatusOK, gin.H{"message": "Menu item updated successfully"})
		})

		menuRoutes.DELETE("/:itemId", func(c *gin.Context) {
			role, _ := c.Get("role")
			if role != "admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
				return
			}
			itemId := c.Param("itemId")
			if err := initializers.DB.Delete(&models.Menu{}, itemId).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete menu item"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted successfully"})
		})

		menuRoutes.GET("/", func(c *gin.Context) {
			var menuItems []models.Menu
			if err := initializers.DB.Find(&menuItems).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve menu items"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"menuItems": menuItems})
		})
	}
}

func setupBasketEndpoints(router *gin.Engine) {
	router.GET("/basket", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "welcome to basket endpoint"})
	})
}

<<<<<<< HEAD
=======
type OrderRequest struct {
	OrderItems []OrderItem `json:"order_items"`
}

type OrderItem struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func setupOrderEndpoints(router *gin.Engine) {
	orders := router.Group("/orders")
	{
		// POST запрос для создания нового заказа
		orders.POST("/", utils.AuthMiddleware(), func(c *gin.Context) {
			// Получаем UserID из контекста, предоставленного middleware
			userIDAny, _ := c.Get("ID")
			userID, ok := userIDAny.(uint)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}

			var orderReq OrderRequest
			if err := c.BindJSON(&orderReq); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			// Создаем новый заказ
			newOrder := models.Order{
				UserID:       userID,
				OrderDetails: make([]models.OrderDetail, 0),
				CreatedAt:    time.Now(),
			}

			// Вычисляем общую стоимость заказа
			var totalPrice decimal.Decimal

			// Добавляем элементы заказа и вычисляем общую стоимость
			for _, item := range orderReq.OrderItems {
				menuItem := models.Menu{}
				result := initializers.DB.First(&menuItem, item.ProductID)
				if result.Error != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
					return
				}

				orderDetail := models.OrderDetail{
					ItemID:   item.ProductID,
					Quantity: item.Quantity,
				}
				newOrder.OrderDetails = append(newOrder.OrderDetails, orderDetail)

				// Вычисляем стоимость для данного элемента заказа и добавляем к общей стоимости
				itemPrice := menuItem.Price.Mul(decimal.NewFromInt(int64(item.Quantity)))
				totalPrice = totalPrice.Add(itemPrice)
			}

			newOrder.TotalPrice = totalPrice

			// Сохраняем новый заказ в базе данных
			if err := initializers.DB.Create(&newOrder).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
				return
			}

			c.JSON(http.StatusCreated, newOrder)
		})
	}
}
>>>>>>> 1821d9010e77f7d3debd2018a1fb15da611554ab
