package productcontroller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"
	"warehouse-management/config"
	"warehouse-management/entities"
	productmodel "warehouse-management/models/productModel"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	product, err := productmodel.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid getting all products"})
		return
	}
	// Mengatur header content type sebagai JSON
	w.Header().Set("Content-Type", "application/json")
	// Encode data kategori ke dalam format JSON dan kirimkan ke client
	err = json.NewEncoder(w).Encode(map[string]any{
		"product": product,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid make a json response"})
		return
	}
}

func GetAllProduct(w http.ResponseWriter, r *http.Request) {
	// Fetch all products with their associated categories
	products, err := productmodel.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid getting all products"})
		return
	}

	// Parse the HTML template
	temp, err := template.ParseFiles("views/products/index.html")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid parse file"})
		return
	}

	// Pass products to the template
	err = temp.Execute(w, map[string]any{
		"products": products,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to render template"})
		return
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close() // Close the body after reading

	// Check if the body is empty
	if len(body) == 0 {
		http.Error(w, "Request body cannot be empty", http.StatusBadRequest)
		return
	}

	// Decode JSON request body into struct Product
	var product entities.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input (e.g., product name is required)
	if product.Name == "" {
		http.Error(w, "Product name is required", http.StatusBadRequest)
		return
	}

	if product.Description == "" {
		http.Error(w, "Product Description is required", http.StatusBadRequest)
		return
	}

	if product.Stock <= 0 {
		http.Error(w, "Product stock is required", http.StatusBadRequest)
		return
	}

	if product.CategoryID <= 0 {
		http.Error(w, "Product category is required", http.StatusBadRequest)
		return
	}

	// Check if the category exists before creating the product
	exists, err := productmodel.CategoryExists(product.CategoryID)
	if err != nil {
		http.Error(w, "Error checking category: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Category does not exist", http.StatusBadRequest)
		return
	}

	// Now create the product
	err = productmodel.Create(&product)
	if err != nil {
		http.Error(w, "Failed to create product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// If successful, send JSON response with product data
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Product added successfully",
		"product": product, // Changed "products" to "product" for singular
	})
}

func AddNewProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Fetch categories from the database
		categories, err := productmodel.GetCategories()
		if err != nil {
			http.Error(w, "Failed to fetch categories: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Render the add product form and pass categories to the template
		tmpl, err := template.ParseFiles("views/products/create.html") // Adjust path as necessary
		if err != nil {
			http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, map[string]interface{}{
			"Categories": categories,
		})
		if err != nil {
			http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "POST" {
		// Decode form values into struct Product
		var product entities.Product
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
			return
		}
		// Debug: Print all form values
		fmt.Println("Form Values:", r.Form)

		product.Name = r.FormValue("nama_product")
		// Validate input (e.g., product name is required)
		if product.Name == "" {
			http.Error(w, "Product name is required", http.StatusBadRequest)
			return
		}
		product.Description = r.FormValue("deskripsi")
		// Convert stock and category_id from string to int
		stockStr := r.FormValue("stok")
		if stockStr == "" {
			http.Error(w, "Stock value is required", http.StatusBadRequest)
			return
		}

		stock, err := strconv.Atoi(stockStr)
		if err != nil {
			http.Error(w, "Invalid stock value: must be a number", http.StatusBadRequest)
			return
		}
		product.Stock = stock

		categoryID, err := strconv.Atoi(r.FormValue("category_id"))
		if err != nil {
			http.Error(w, "Invalid category ID: "+err.Error(), http.StatusBadRequest)
			return
		}
		product.CategoryID = uint(categoryID)

		// Check if the category exists before creating the product
		exists, err := productmodel.CategoryExists(product.CategoryID)
		if err != nil {
			http.Error(w, "Error checking category: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "Category does not exist", http.StatusBadRequest)
			return
		}

		// Now create the product
		err = productmodel.Create(&product)
		if err != nil {
			http.Error(w, "Failed to create product: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// If successful, redirect or send a success message
		http.Redirect(w, r, "/product", http.StatusSeeOther)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productIDStr := vars["id"]
	// Konversi categoryIDStr ke uint
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Cek metode HTTP
	if r.Method == "GET" {
		// Ambil data kategori dari database berdasarkan ID
		product, err := productmodel.GetByID(uint(productID))
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		// Kirim data kategori dalam format JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
		return
	}

	if r.Method == "PUT" || r.Method == "PATCH" {
		var updateProduct entities.Product
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&updateProduct); err != nil {
			http.Error(w, "Invalid category data", http.StatusBadRequest)
			return
		}

		// Periksa apakah produk ada sebelum memperbarui
		existingProduct, err := productmodel.GetByID(uint(productID))
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		// Update kategori di database
		updateProduct.ID = existingProduct.ID // Convert ke uint setelah parse uint64
		err = productmodel.Update(&updateProduct)
		if err != nil {
			http.Error(w, "Failed to update category", http.StatusInternalServerError)
			return
		}

		// Kembalikan response sukses
		w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content
		return
	}
}

func EditNewProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)    // Get the URL parameters
	idString := vars["id"] // Extract the "id" parameter from the URL

	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/products/edit.html")
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to parse template"})
			return
		}

		if idString == "" {
			json.NewEncoder(w).Encode(map[string]string{"error": "missing id parameter"})
			return
		}

		IDint, err := strconv.Atoi(idString)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid id parameter"})
			return
		}

		product := &entities.Product{}
		err = config.DB.Preload("Category").First(product, IDint).Error
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "product not found"})
			return
		}

		// Fetch categories for the dropdown
		var categories []entities.Category
		err = config.DB.Find(&categories).Error
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to load categories"})
			return
		}

		// Prepare data for the template
		data := map[string]any{
			"product":    product,
			"categories": categories, // Passing categories to the template
		}

		err = temp.Execute(w, data)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to execute template"})
			return
		}
		return
	}

	if r.Method == "POST" {
		if idString == "" {
			json.NewEncoder(w).Encode(map[string]string{"error": "missing id parameter"})
			return
		}

		IDint, err := strconv.Atoi(idString)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid id parameter"})
			return
		}

		// Get updated product data from the form
		name := r.FormValue("name")
		description := r.FormValue("description")
		stock, err := strconv.Atoi(r.FormValue("stock"))
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid stock"})
			return
		}
		categoryID, err := strconv.Atoi(r.FormValue("category_id"))
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid category"})
			return
		}

		// Find the existing product
		product := &entities.Product{}
		err = config.DB.First(product, IDint).Error
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "product not found"})
			return
		}

		// Update product details
		product.Name = name
		product.Description = description
		product.Stock = stock
		product.CategoryID = uint(categoryID)

		// Save the updated product
		err = config.DB.Save(product).Error
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to update product"})
			return
		}

		// Redirect to the product list page after successful update
		http.Redirect(w, r, "/product", http.StatusFound)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing id parameter"})
		return
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id parameter"})
	}

	err = productmodel.Delete(idInt)
	if err != nil {
		if err.Error() == "product not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "product not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete product"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "product deleted successfully"})
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	if idString == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Missing id parameter"})
		return
	}
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid id parameter"})
		return
	}
	err = productmodel.Delete(idInt)
	if err != nil {
		http.Error(w, "Failed to delete product: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/product", http.StatusSeeOther)
}
