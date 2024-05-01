package main

import (
	"final_project/initializers"
	"github.com/gin-gonic/gin"
	"final_project/internal/router"
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
	router := router.SetupRouter()

}
