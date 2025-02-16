package customercontoller

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"warehouse-management/entities"
	customermodel "warehouse-management/models/customerModel"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	customer, err := customermodel.GetAllCustomer()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid getting all customers"})
	}
	// Mengatur header content type sebagai JSON
	w.Header().Set("Content-Type", "application/json")
	// Encode data kategori ke dalam format JSON dan kirimkan ke client
	err = json.NewEncoder(w).Encode(map[string]any{
		"customer": customer,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid make a json response"})
		return
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if r.Body == nil {
		http.Error(w, "Request body cannot be empty", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var customer entities.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validation
	if customer.Name == "" {
		http.Error(w, "Customer name is required", http.StatusBadRequest)
		return
	}

	if customer.Address == "" {
		http.Error(w, "Customer address is required", http.StatusBadRequest)
		return
	}

	if customer.Email == "" {
		http.Error(w, "Customer email is required", http.StatusBadRequest)
		return
	}

	// Validate email format
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if match, _ := regexp.MatchString(emailRegex, customer.Email); !match {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if customer.Phone == "" {
		http.Error(w, "Customer phone is required", http.StatusBadRequest)
		return
	}

	// Validate phone number (only digits)
	phoneRegex := `^[0-9]+$`
	if match, _ := regexp.MatchString(phoneRegex, customer.Phone); !match {
		http.Error(w, "Invalid phone number format", http.StatusBadRequest)
		return
	}

	// Create customer in the database
	err = customermodel.CreateCustomer(&customer)
	if err != nil {
		http.Error(w, "Failed to create customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message":  "Customer added successfully",
		"customer": customer,
	})
}

func Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerIDStr := vars["id"]

	// Convert customerIDStr to uint64
	customerID, err := strconv.ParseUint(customerIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Fetch customer details for GET request
	if r.Method == http.MethodGet {
		customer, err := customermodel.GetByID(uint(customerID))
		if err != nil {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customer)
		return
	}

	// Handle PUT or PATCH for updating customer data
	if r.Method == http.MethodPut || r.Method == http.MethodPatch {
		if r.Body == nil {
			http.Error(w, "Request body cannot be empty", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Fetch existing customer to avoid unintended overwrites
		existingCustomer, err := customermodel.GetByID(uint(customerID))
		if err != nil {
			http.Error(w, "Customer not found", http.StatusNotFound)
			return
		}

		// Decode JSON request body
		var updateCustomer entities.Customer
		if err := json.NewDecoder(r.Body).Decode(&updateCustomer); err != nil {
			http.Error(w, "Invalid customer data", http.StatusBadRequest)
			return
		}

		// Ensure the update doesn't erase existing values if fields are empty
		if updateCustomer.Name != "" {
			existingCustomer.Name = updateCustomer.Name
		}
		if updateCustomer.Phone != "" {
			existingCustomer.Phone = updateCustomer.Phone
		}
		if updateCustomer.Address != "" {
			existingCustomer.Address = updateCustomer.Address
		}
		if updateCustomer.Email != "" {
			existingCustomer.Email = updateCustomer.Email
		}

		// Update customer in the database
		err = customermodel.Update(&existingCustomer)
		if err != nil {
			http.Error(w, "Failed to update customer", http.StatusInternalServerError)
			return
		}

		// Respond with updated customer data
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"message":  "Customer updated successfully",
			"customer": existingCustomer,
		})
		return
	}

	// If method is not supported
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	err = customermodel.Delete(idInt)
	if err != nil {
		if err.Error() == "customer not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "customer not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete customer"})
		return
	}

	// Return success response with no content
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "customer deleted successfully"})
}
