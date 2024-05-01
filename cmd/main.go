package main

import (
	"final_project/initializers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func init() {
	initializers.GetKeysInEnv()
	initializers.ConnectDb()
}

// @title						Gallery service
//
//	@version					1.0
//	@description				Gallery service
//	@termsOfService				http://swagger.io/terms/
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						localhost:8082
//	@securityDefinitions.apikey	BearerAuth
//	@type						apiKey
//	@name						Authorization
//	@in							header
//	@schemes
func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}

func PingExample(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
