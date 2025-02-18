package categoriescontroller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"
	"warehouse-management/config"
	"warehouse-management/entities"
	categoriesmodel "warehouse-management/models/categoriesModel"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	categories, err := categoriesmodel.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entities.NewErrorResponse("Failed to get categories"))
		return
	}
	// Mengatur header content type sebagai JSON
	w.Header().Set("Content-Type", "application/json")
	// Encode data kategori ke dalam format JSON dan kirimkan ke client
	err = json.NewEncoder(w).Encode(map[string]any{
		"categories": categories,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entities.NewErrorResponse("Invalid make a json response"))
		return
	}
}

func GetCategoriesAll(w http.ResponseWriter, r *http.Request) {
	categories, err := categoriesmodel.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entities.NewErrorResponse("Failed to get all categories"))
		return
	}
	tmpl, err := template.ParseFiles("views/category/index.html")
	if err != nil {
		http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, map[string]any{
		"categories": categories,
	})
	if err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	// Hanya izinkan metode POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entities.NewErrorResponse("MethodNotAllowed"))
		return
	}

	// Periksa jika body kosong
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entities.NewErrorResponse("Body cannot be empty"))
		return
	}

	// Dekode JSON request body ke struct Category
	var category entities.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entities.NewErrorResponse("Invalid request body"))
		return
	}

	// Validasi input (misal nama kategori wajib diisi)
	if category.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entities.NewErrorResponse("Category name is empty"))
		return
	}

	// Panggil fungsi model untuk menyimpan kategori baru
	err = categoriesmodel.Create(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entities.NewErrorResponse("Failed to create category"))
		return
	}

	// Jika berhasil, kirim respons JSON dengan data kategori
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message":  "Category added successfully",
		"category": category,
	})
}

func AddNewCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/category/create.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entities.NewErrorResponse("Category name is empty"))
			return
		}
		temp.Execute(w, nil)
	}

	if r.Method == "POST" {
		var category entities.Category

		category.Name = r.FormValue("name")
		category.CreatedAt = time.Now()
		category.UpdatedAt = time.Now()

		ok := categoriesmodel.Create(&category)
		if ok != nil {
			temp, _ := template.ParseFiles("views/category/create.html")
			temp.Execute(w, nil)
		}

		http.Redirect(w, r, "/categories", http.StatusSeeOther)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari parameter URL
	vars := mux.Vars(r)
	categoryIDStr := vars["id"]

	// Konversi categoryIDStr ke uint
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Cek metode HTTP
	if r.Method == "GET" {
		// Ambil data kategori dari database berdasarkan ID
		category, err := categoriesmodel.GetByID(uint(categoryID))
		if err != nil {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		// Kirim data kategori dalam format JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(category)
		return
	}

	// Jika menggunakan PUT atau PATCH
	if r.Method == "PUT" || r.Method == "PATCH" {
		// Ambil data kategori yang dikirimkan oleh client
		var updatedCategory entities.Category
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&updatedCategory); err != nil {
			http.Error(w, "Invalid category data", http.StatusBadRequest)
			return
		}

		excitingCategory, err := categoriesmodel.GetByID(uint(categoryID))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entities.NewErrorResponse("Category Not Found"))
			return
		}

		// Update kategori di database
		updatedCategory.ID = excitingCategory.ID // Convert ke uint setelah parse uint64
		err = categoriesmodel.Update(&updatedCategory)
		if err != nil {
			http.Error(w, "Failed to update category", http.StatusInternalServerError)
			return
		}

		// Kembalikan response sukses
		w.WriteHeader(http.StatusNoContent) // HTTP 204 No Content
		return
	}
}

func EditNewCategories(w http.ResponseWriter, r *http.Request) {
	// Handle GET request (display the edit form)
	if r.Method == "GET" {
		// Parse the template
		temp, err := template.ParseFiles("views/category/edit.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entities.NewErrorResponse("Failed to parse template"))
			return
		}

		// Get the "id" query parameter
		idString := r.URL.Query().Get("id")
		if idString == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entities.NewErrorResponse("Missing id parameter"))
			return
		}

		// Convert the "id" parameter to an integer
		idInt, err := strconv.Atoi(idString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entities.NewErrorResponse("Invalid id parameter"))
			return
		}

		// Fetch the category data from the database
		category := &entities.Category{}
		err = config.DB.First(category, idInt).Error
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(entities.NewErrorResponse("Category not found"))
			return
		}

		// Prepare data for the template
		data := map[string]any{
			"category": category,
		}

		// Execute the template
		err = temp.Execute(w, data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entities.NewErrorResponse("Failed to execute template"))
			return
		}
		return
	}

	// Handle POST request (update the category)
	if r.Method == "POST" {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entities.NewErrorResponse("Failed to parse form data"))
			return
		}

		// Get the "id" parameter from the form
		idString := r.FormValue("id")
		if idString == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entities.NewErrorResponse("Missing id parameter"))
			return
		}

		// Convert the "id" parameter to an integer
		idInt, err := strconv.Atoi(idString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(entities.NewErrorResponse("Invalid id parameter"))
			return
		}

		// Fetch the existing category from the database
		category := &entities.Category{}
		err = config.DB.First(category, idInt).Error
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(entities.NewErrorResponse("Category not found"))
			return
		}

		// Update the category fields from the form data
		category.Name = r.FormValue("name") // Example: Update the name field

		err = categoriesmodel.Update(category)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to update category"})
			return
		}

		http.Redirect(w, r, "/categories", http.StatusSeeOther)
		return
	}

	// Handle unsupported methods
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	// Extract the ID from the URL path
	vars := mux.Vars(r)
	idString, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing id parameter"})
		return
	}

	// Convert the ID to an integer
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id parameter"})
		return
	}

	// Call the model's Delete function
	err = categoriesmodel.Delete(idInt)
	if err != nil {
		if err.Error() == "category not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(entities.NewErrorResponse("Category not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete category"})
		return
	}

	// Return success response with no content
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "category deleted successfully"})
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}
	err = categoriesmodel.Delete(idInt)
	if err != nil {
		http.Error(w, "Failed to delete category: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}
