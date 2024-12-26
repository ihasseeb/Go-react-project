package main

// Necessary imports
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers" // Import handlers package
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// User struct to handle incoming JSON data
type User struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone"`
	DOB         string `json:"dob"`
	Address     string `json:"address"`
}

// Response struct for the outgoing JSON response
type Response struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

// saveDataToFile saves the response data to a JSON file, appending to the file
func saveDataToFile(response Response, fileName string) error {
	// Open the file, or create it if it doesn't exist
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Marshal the new response data into JSON
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	// Add a newline after each JSON object to separate them
	data = append(data, []byte("\n")...)

	// Write the new data to the file (append mode)
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing data to file: %v", err)
	}

	return nil
}

// this will initialize the database
func initDB() (*sql.DB, error) {
	// please insert your own DB_USER_NAME and PASSWORD and DB_NAME
	connStr := "user=user_name password=passord dbname=database_name sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	return db, nil
}

// this will create a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// validating the DOB
	_, err = time.Parse("2006-01-02", user.DOB)
	if err != nil {
		http.Error(w, "Invalid Date of Birth format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// connect database
	db, err := initDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("Database connection failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// send the query
	query := `INSERT INTO users (first_name, last_name, email, phone_number, dob, address)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var userID int
	err = db.QueryRow(query, user.FirstName, user.LastName, user.Email, user.PhoneNumber, user.DOB, user.Address).Scan(&userID)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		http.Error(w, fmt.Sprintf("Failed to insert user: %v", err), http.StatusInternalServerError)
		return
	}

	// Create the response data
	response := Response{
		ID:       userID,
		FullName: user.FirstName + " " + user.LastName,
		Email:    user.Email,
	}

	// Save the data to a JSON file
	err = saveDataToFile(response, "users.json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error saving data: %v", err), http.StatusInternalServerError)
		return
	}

	// Send the response back to the client
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// main function
func main() {
	// setting the router
	r := mux.NewRouter()

	r.HandleFunc("/api/users", createUser).Methods("POST")

	// Enable CORS
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	// starting the server with CORS middleware
	fmt.Println("Server started at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
