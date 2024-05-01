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
	setupOrderRoutes(router)

	//only admin routers
	SetupOrderUpdateRouter(router)
	return router
}

func setupPublicEndpoints(router *gin.Engine) {
	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "welcome to public endpoint"})
	})
}

func setupPrivateEndpoints(router *gin.Engine) {
	router.GET("/private", utils.AuthMiddleware(), func(c *gin.Context) {
		role, _ := c.Get("role")
		if role == "admin" {
			c.JSON(http.StatusOK, gin.H{"message": "welcome to private endpoint (menu admin)"})
		} else if role == "client" {

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

type OrderRequest struct {
	OrderItems []OrderItem `json:"order_items"`
}

type OrderItem struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func setupOrderEndpoints(router *gin.Engine) {
	orders := router.Group("/orders", utils.AuthMiddleware())
	{
		orders.POST("/", func(c *gin.Context) {
			userID, _ := c.Get("ID")
			var orderReq OrderRequest
			if err := c.BindJSON(&orderReq); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			newOrder := models.Order{
				UserID:       userID.(uint),
				OrderDetails: []models.OrderDetail{},
				CreatedAt:    time.Now(),
				OrderStatus:  models.Preparing,
			}

			var totalPrice decimal.Decimal
			tx := initializers.DB.Begin()

			for _, item := range orderReq.OrderItems {
				menuItem := models.Menu{}
				result := tx.First(&menuItem, item.ProductID)
				if result.Error != nil {
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found", "productID": item.ProductID})
					return
				}

				if menuItem.Quantity < item.Quantity {
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock", "productID": item.ProductID})
					return
				}

				menuItem.Quantity -= item.Quantity
				if err := tx.Save(&menuItem).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu item stock", "productID": item.ProductID})
					return
				}

				itemTotalCost := menuItem.Price.Mul(decimal.NewFromInt(int64(item.Quantity)))

				orderDetail := models.OrderDetail{
					ItemID:    item.ProductID,
					Quantity:  item.Quantity,
					TotalCost: itemTotalCost,
				}
				newOrder.OrderDetails = append(newOrder.OrderDetails, orderDetail)

				totalPrice = totalPrice.Add(itemTotalCost)
			}

			newOrder.TotalPrice = totalPrice

			if err := tx.Create(&newOrder).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
				return
			}

			tx.Commit()
			c.JSON(http.StatusCreated, newOrder)
		})
	}
}

func setupOrderRoutes(router *gin.Engine) {
	orders := router.Group("/orders", utils.AuthMiddleware())
	{
		orders.GET("/", func(c *gin.Context) {
			userID, exists := c.Get("ID")
			if !exists {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found"})
				return
			}

			var userOrders []models.Order
			if err := initializers.DB.Preload("OrderDetails.MenuItem").Where("user_id = ?", userID.(uint)).Find(&userOrders).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders", "details": err.Error()})
				return
			}

			response := make([]map[string]interface{}, 0)
			for _, order := range userOrders {
				orderItems := make([]map[string]interface{}, 0)
				for _, detail := range order.OrderDetails {
					orderItems = append(orderItems, map[string]interface{}{
						"id": detail.ID,
						"item": map[string]interface{}{
							"ID":          detail.MenuItem.ID,
							"name":        detail.MenuItem.Name,
							"description": detail.MenuItem.Description,
							"price":       detail.MenuItem.Price.String(),
						},
						"quantity":    detail.Quantity,
						"total_price": detail.TotalCost.String(),
					})
				}

				response = append(response, map[string]interface{}{

					"order_id":     order.ID,
					"order_items":  orderItems,
					"order_status": order.OrderStatus,
					"order_cost":   order.TotalPrice.String(),
					"updated_at":   order.UpdatedAt.Format(time.RFC3339Nano),
					"created_at":   order.CreatedAt.Format(time.RFC3339Nano),
				})
			}

			c.JSON(http.StatusOK, response)
		})
	}
}

type UpdateOrderData struct {
	Status string `json:"status" binding:"required"`
}

func SetupOrderUpdateRouter(router *gin.Engine) {
	orders := router.Group("/orders", utils.AuthMiddleware())
	{
		orders.PATCH("/:OrderId", func(c *gin.Context) {
			var updateData UpdateOrderData
			if err := c.BindJSON(&updateData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
				return
			}

			role, _ := c.Get("role")
			if role != "admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can update order status"})
				return
			}

			orderID := c.Param("OrderId")
			orderStatus := models.Status(updateData.Status)

			switch orderStatus {
			case models.Canceled, models.Preparing, models.Ready, models.Completed:
				order := &models.Order{}
				result := initializers.DB.First(order, orderID)
				if result.Error != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
					return
				}

				order.OrderStatus = orderStatus

				// Сохраняем изменения и выполняем проверку перед сохранением
				if err := initializers.DB.Save(order).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status", "details": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status"})
			}
		})
	}
}
