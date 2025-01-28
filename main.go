package main

import (
	"log"
	"net/http"
	"warehouse-management/config"
	categoriescontroller "warehouse-management/controller/categoriesController"
	homecontroller "warehouse-management/controller/homeController"
	productcontroller "warehouse-management/controller/productController"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()
	r := mux.NewRouter()

	r.HandleFunc("/", homecontroller.Welcome)

	//1. Category API
	r.HandleFunc("/api/categories", categoriescontroller.Index).Methods("GET")
	r.HandleFunc("/api/categories/add", categoriescontroller.Add).Methods("POST")
	r.HandleFunc("/api/categories/edit/{id}", categoriescontroller.Edit).Methods("GET", "PUT", "PATCH")
	r.HandleFunc("/api/categories/delete", categoriescontroller.Delete).Methods("DELETE")

	//2. Category Views
	r.HandleFunc("/categories", categoriescontroller.GetCategoriesAll).Methods("GET")
	r.HandleFunc("/categories/add", categoriescontroller.AddNewCategories).Methods("GET")  // Menampilkan form untuk menambahkan kategori
	r.HandleFunc("/categories/add", categoriescontroller.AddNewCategories).Methods("POST") // melakukan eksekusi untuk menambahkan kategori
	r.HandleFunc("/categories/edit", categoriescontroller.EditNewCategories).Methods("GET", "POST")
	r.HandleFunc("/categories/delete", categoriescontroller.Delete).Methods("DELETE")

	//3. Products
	r.HandleFunc("/api/product", productcontroller.Index).Methods("GET")
	r.HandleFunc("/api/product/add", productcontroller.Add).Methods("POST")
	r.HandleFunc("/api/product/update", productcontroller.Update).Methods("POST")
	r.HandleFunc("/api/product/detail", productcontroller.Detail).Methods("GET")
	r.HandleFunc("/api/product/delete", productcontroller.Delete).Methods("DELETE")

	log.Println("server running on port 8080")
	http.ListenAndServe(":8080", r)
}
