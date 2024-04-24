package main

import (
	"final_project/initializers"
)

func init() {
	initializers.GetKeysInEnv()
	initializers.ConnectDb()
}

func main() {

}
