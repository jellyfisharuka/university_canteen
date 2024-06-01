package main

import (
	"final_project/initializers"
	"final_project/internal/router"
	"log"
)

func init() {
	initializers.GetKeysInEnv()
	initializers.ConnectDb()
}

func main() {
	
	router := router.SetupRouter()
	err := router.Run(":8092")
if err != nil {
    log.Fatal(err)
}

}
