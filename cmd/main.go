package main

import (
	"context"
	"final_project/initializers"
	"final_project/internal/router"
	"fmt"
	"log"

	"github.com/jomei/notionapi"
)

func init() {
	initializers.GetKeysInEnv()
	initializers.ConnectDb()
}

func main() {
	client := notionapi.NewClient("secret_6Z366pdbdqmV5H9RP5K2ksjoA1W44E6pUR6A5YSnguq")
	page, err := client.Page.Get(context.Background(), "your_page_id")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Page ID:", page.ID)
    //fmt.Println("Page Title:", page.Title)
    fmt.Println("Page URL:", page.URL)
	router := router.SetupRouter()
	error := router.Run(":8092")
if error != nil {
    log.Fatal(err)
}

}
