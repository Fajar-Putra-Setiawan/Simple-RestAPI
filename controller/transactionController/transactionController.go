package transactioncontroller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"warehouse-management/entities"
	transactionmodel "warehouse-management/models/transactionModel"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	transaction, err := transactionmodel.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid getting all transaction"})
		return
	}
	err = json.NewEncoder(w).Encode(map[string]any{
		"transaction": transaction,
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		http.Error(w, "Request Body cant be empty", http.StatusBadRequest)
	}

	var transaction entities.Transaction
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validasi data
	if transaction.ProductID == 0 || transaction.Quantity <= 0 || transaction.TransactionType == "" {
		http.Error(w, "Invalid transaction data", http.StatusBadRequest)
		return
	}

	// Hitung total amount
	transaction.TotalAmount = float64(transaction.Quantity) * transaction.Price

	// Simpan transaksi dan update stok
	if err := transactionmodel.CreateTransaction(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Kirim response sukses
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message":     "Transaction added successfully",
		"transaction": transaction,
	})
}

func Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionIDstr := vars["id"]

	transactionID, err := strconv.ParseUint(transactionIDstr, 10, 16)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodGet {
		// Ambil data produk dari database berdasarkan ID
		transaction, err := transactionmodel.GetByID(uint(transactionID))
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		// Kirim data produk dalam format JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transaction)
		return
	}
}
