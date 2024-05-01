package router

import (
	"final_project/initializers"
	"final_project/internal/models"
	"final_project/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "welcome to public endpoint"})
	})
	router.GET("/private", utils.AuthMiddleware(), func(c *gin.Context) {
		role, _ := c.Get("role") // Получаем роль пользователя из контекста
		if role == "menu_admin" {
			// Действия для администратора меню
			c.JSON(http.StatusOK, gin.H{"message": "welcome to private endpoint (menu admin)"})
		} else {
			// Действия для обычного пользователя
			c.JSON(http.StatusOK, gin.H{"message": "welcome to private endpoint (user)"})
		}
	})
	router.POST("/login", func(c *gin.Context) {
		var loginUser models.User
		if err := c.BindJSON(&loginUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		var existingUser models.User
		result := initializers.DB.Where("username = ?", loginUser.Username).First(&existingUser)
		if result.Error != nil || !utils.CheckPassword(existingUser.Password, loginUser.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		token, err := utils.GenerateToken(loginUser.Username)
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

	return router
}