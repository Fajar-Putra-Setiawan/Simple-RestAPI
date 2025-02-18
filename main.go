package main

import (
	"log"
	"net/http"
	"warehouse-management/api"
	"warehouse-management/config"
)

func main() {
	config.ConnectDB()
	r := api.NewRouter()

	log.Println("server running on port 8080")
	http.ListenAndServe(":8080", r)
}
