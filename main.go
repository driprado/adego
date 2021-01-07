// @TODO:
// 1 - Implement SQL DB into it
	//a) https://www.youtube.com/watch?v=DWNozbk_fuk
	// b) https://www.youtube.com/watch?v=pr4x8KdIfDU <-- FOCUS HERE
// 2 - Make it architecture with go-kit
// 3 - Plug into DynamoDB and APIGateway

// @NOTE:
// github.com/driprado/starting_up_with_go/go_nethttp/restapi for more details on golang restAPIS

// To test this code run locally:
// Install VS Code Remote container plugin
// Use in `Reopen in container` mode
// Set local github credentials
// go build
// ./adega
// It's served at: http://localhost:8000/api/wines
// Test routes with postman

package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Wine struct, holds wines attributes
type Wine struct {
	ID            string `json:"id"`
	Country       string `json:"country"`
	Region        string `json:"region"`
	Name          string `json:"name"`
	Producer      string `json:"producer"`
	Grape         string `json:"grape"`
	Vintage       string `json:"vintage"`
	PurchasePrice string `json:"purchaseprice"`
	CurrentPrice  string `json:"currentprice"`
	PurchaseDate  string `json:"purchasedate"`
	Stars         *Stars `json:"starts"`
}

// Stars struct holds wine subattribute
type Stars struct {
	Personal  string `json:"personal"`
	Community string `json:"community"`
}

// Declare wines slice of type Wine
var wines []Wine

// json
// ENCODE Response from HTTP
// DECODE Request from HTTP only on POST and PUT

// getWines returns  all wines (GET)
func getWines(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Set ResponseWriter header to json
	json.NewEncoder(w).Encode(wines)                   // ResponseWrier to encode slice of wines into json
}

// createWine function
func createWine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newWine Wine
	_ = json.NewDecoder(r.Body).Decode(&newWine)  // Decode request
	newWine.ID = strconv.Itoa(rand.Intn(1000000)) // Generate random ID number and convert to string
	wines = append(wines, newWine)
	json.NewEncoder(w).Encode(newWine) // Encode response

}

// Get single wine function Grab by what? everything, grape, year, price range
func getWine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Get parameters from URL
	parameters := mux.Vars(r) //func Vars(r *http.Request) map[string]string | returns the url variables for the current request.
	// Loop through wines until find ID
	for _, wine := range wines {
		if wine.ID == parameters["id"] { // "id" is defined in the route, ex: "/api/wines/{id}"
			json.NewEncoder(w).Encode(wine)
			return
		}
	}
	json.NewEncoder(w).Encode(&Wine{}) // Encode everything in the Wine structure

}

// updateWine
func updateWine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := mux.Vars(r)
	for index, item := range wines {
		if item.ID == parameters["id"] {
			wines = append(wines[:index], wines[index+1:]...)
			var wine Wine
			_ = json.NewDecoder(r.Body).Decode(&wine)
			wine.ID = parameters["id"]
			wines = append(wines, wine)
			json.NewEncoder(w).Encode(wine)
			return
		}
	}
}

// deleteFromSlice
func deleteFromSlice(s []Wine, index int) []Wine {
	return append(s[:index], s[index+1:]...)
}

func deleteWine(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	parameters := mux.Vars(r)
	for i, wine := range wines {
		if wine.ID == parameters["id"] {
			wines = deleteFromSlice(wines, i)
			break
		}
	}
	json.NewEncoder(w).Encode(wines)
}

func main() {
	//Initiare mux router
	router := mux.NewRouter()

	wines = append(wines, Wine{ID: "1", Country: "australia", Region: "Coonawarra", Name: "prosperitas", Producer: "prosperitas", Grape: "shiraz", Vintage: "2016", PurchasePrice: "80", CurrentPrice: "160", PurchaseDate: "2018", Stars: &Stars{Personal: "5", Community: "1"}})

	wines = append(wines, Wine{ID: "2", Country: "france", Region: "rhone", Name: "Guigal Cotes du Rhone", Producer: "Guigal", Grape: "Grenache Shiraz Mourvedre", Vintage: "2016", PurchasePrice: "12", CurrentPrice: "20", PurchaseDate: "2016", Stars: &Stars{Personal: "3", Community: "5"}})

	// Route Handlers / Endpoints
	// func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) // pattern means uri path in this case
	// HandleFunc registers the handler function for the given pattern in the DefaultServeMux.
	router.HandleFunc("/api/wines", getWines).Methods("GET")           // Get all wines
	router.HandleFunc("/api/wines/{id}", getWine).Methods("GET")       // Get single wine by id
	router.HandleFunc("/api/wines", createWine).Methods("POST")        // Create wine
	router.HandleFunc("/api/wines/{id}", updateWine).Methods("PUT")    // Update wine
	router.HandleFunc("/api/wines/{id}", deleteWine).Methods("DELETE") // Delete wine

	// Serve function
	log.Fatal(http.ListenAndServe(":8000", router)) // wrapped in log.Fatal to log in case of error

}
