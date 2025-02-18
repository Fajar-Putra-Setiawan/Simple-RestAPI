package api

import (
	categoriescontroller "warehouse-management/controller/categoriesController"
	customercontoller "warehouse-management/controller/customerContoller"
	homecontroller "warehouse-management/controller/homeController"
	productcontroller "warehouse-management/controller/productController"
	transactioncontroller "warehouse-management/controller/transactionController"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", homecontroller.Welcome)

	//1. Category API
	r.HandleFunc("/api/categories", categoriescontroller.Index).Methods("GET")
	r.HandleFunc("/api/categories/add", categoriescontroller.Add).Methods("POST")
	r.HandleFunc("/api/categories/edit/{id}", categoriescontroller.Edit).Methods("GET", "PUT", "PATCH")
	r.HandleFunc("/api/categories/delete/{id}", categoriescontroller.Delete).Methods("DELETE")

	//2. Products API
	r.HandleFunc("/api/product", productcontroller.Index).Methods("GET")
	r.HandleFunc("/api/product/add", productcontroller.Add).Methods("POST")
	r.HandleFunc("/api/product/edit/{id}", productcontroller.Edit).Methods("GET", "PUT", "PATCH")
	r.HandleFunc("/api/product/delete/{id}", productcontroller.Delete).Methods("DELETE")

	//3. Category API
	r.HandleFunc("/api/customer", customercontoller.Index).Methods("GET")
	r.HandleFunc("/api/customer/add", customercontoller.Add).Methods("POST")
	r.HandleFunc("/api/customer/edit/{id}", customercontoller.Edit).Methods("GET", "PUT", "PATCH")
	r.HandleFunc("/api/customer/delete/{id}", customercontoller.Delete).Methods("DELETE")

	//4. Transaction API
	r.HandleFunc("/api/transaction", transactioncontroller.Index).Methods("GET")
	r.HandleFunc("/api/transaction/add", transactioncontroller.Add).Methods("POST")
	r.HandleFunc("/api/transaction/edit/{id}", customercontoller.Edit).Methods("GET", "PUT", "PATCH")
	r.HandleFunc("/api/transaction/delete/{id}", customercontoller.Delete).Methods("DELETE")

	//5. Category Views
	r.HandleFunc("/categories", categoriescontroller.GetCategoriesAll).Methods("GET")
	r.HandleFunc("/categories/add", categoriescontroller.AddNewCategories).Methods("GET")  // Menampilkan form untuk menambahkan kategori
	r.HandleFunc("/categories/add", categoriescontroller.AddNewCategories).Methods("POST") // melakukan eksekusi untuk menambahkan kategori
	r.HandleFunc("/categories/edit", categoriescontroller.EditNewCategories).Methods("GET", "POST")
	r.HandleFunc("/categories/delete", categoriescontroller.DeleteCategory).Methods("DELETE", "GET")

	//6. Product Views
	r.HandleFunc("/product", productcontroller.GetAllProduct).Methods("GET")
	r.HandleFunc("/product/add", productcontroller.AddNewProduct).Methods("GET")                // Menampilkan form untuk menambahkan kategori
	r.HandleFunc("/product/add", productcontroller.AddNewProduct).Methods("POST")               // melakukan eksekusi untuk menambahkan kategori
	r.HandleFunc("/product/edit/{id}", productcontroller.EditNewProduct).Methods("GET", "POST") // Updated route to include ID
	r.HandleFunc("/product/delete", productcontroller.DeleteProduct).Methods("DELETE", "GET")

	return r
}
