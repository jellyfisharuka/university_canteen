package main

import (
	"final_project/initializers"
	"final_project/internal/router"
)

func init() {
	initializers.GetKeysInEnv()
	initializers.ConnectDb()
}

func main() {
	router := router.SetupRouter()

    // Запуск сервера на порту 8080
    router.Run(":8080")
}
