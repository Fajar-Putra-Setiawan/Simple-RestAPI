# Simple WMS with REST API  

## Description  
**Simple WMS (Warehouse Management System) with REST API** is a lightweight and efficient warehouse management system built using Golang. It provides essential features to manage products, categories, customers, and transactions efficiently.  

## Features  
- **Product Management**: Add, update, delete, and retrieve product details.  
- **Category Management**: Organize products into categories.  
- **Customer Management**: Store and manage customer information.  
- **Transaction Handling**: Record and track transactions within the system.  

## Installation  

1. **Clone the repository**  
   ```bash
   git clone <repository-url>
   cd simple-RestAPI
   ```  

2. **Install dependencies**  
   ```bash
   go mod tidy
   ```  

3. **Set up the database**  
   - Create a database and update the `.env` file with your database credentials.  

4. **Run the application**  
   ```bash
   go run main.go
   ```  

## API Endpoints  
| Method | Endpoint                | Description                |
|--------|-------------------------|----------------------------|
| GET    | `/`                     | Welcome message            |
| GET    | `/api/categories`        | Get all categories         |
| POST   | `/api/categories/add`    | Add a new category         |
| GET, PUT, PATCH | `/api/categories/edit/{id}` | Edit a category |
| DELETE | `/api/categories/delete/{id}` | Delete a category |
| GET    | `/api/product`           | Get all products           |
| POST   | `/api/product/add`       | Add a new product         |
| GET, PUT, PATCH | `/api/product/edit/{id}` | Edit a product |
| DELETE | `/api/product/delete/{id}` | Delete a product |
| GET    | `/api/customer`          | Get all customers         |
| POST   | `/api/customer/add`      | Add a new customer       |
| GET, PUT, PATCH | `/api/customer/edit/{id}` | Edit a customer |
| DELETE | `/api/customer/delete/{id}` | Delete a customer |
| GET    | `/api/transaction`       | Get all transactions      |
| POST   | `/api/transaction/add`   | Add a new transaction    |
| GET, PUT, PATCH | `/api/transaction/edit/{id}` | Edit a transaction |
| DELETE | `/api/transaction/delete/{id}` | Delete a transaction |

## Author  
Created by **Fajar Putra**  
