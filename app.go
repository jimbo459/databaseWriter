package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialise(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initialiseRoutes()
}

func (a *App) initialiseRoutes() {
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8010", a.Router))
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		errorBody := map[string]string{"error": "Invalid product ID"}
		response, _ := json.Marshal(errorBody)
		w.Header().Set("Content-Type", "application-json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}

	product := product{ID: id}

	if err := product.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			errorMessage := map[string]string{"error": "Product not found"}
			body, _ := json.Marshal(errorMessage)

			w.Header().Set("Content-Type", "application-json")
			w.WriteHeader(http.StatusNotFound)
			w.Write(body)

		default:
			errorMessage := map[string]string{"error": err.Error()}
			body, _ := json.Marshal(errorMessage)

			w.Header().Set("Content-Type", "application-json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(body)
		}
		return
	}

	payload, _ := json.Marshal(product)
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {

	response, _ := json.Marshal([]string{})
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var product product

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		errorMessage := map[string]string{"error": "Invalid Payload"}
		body, _ := json.Marshal(errorMessage)

		w.Header().Set("Content-Type", "application-json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
	}

	if err := product.createProduct(a.DB); err != nil {
		errorMessage := map[string]string{"error": err.Error()}
		body, _ := json.Marshal(errorMessage)

		w.Header().Set("Content-Type", "application-json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
	}

	createdProduct, _ := json.Marshal(product)
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusCreated)
	w.Write(createdProduct)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	var product product
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		errorMessage := map[string]string{"error": "Invalid product ID"}
		body, _ := json.Marshal(errorMessage)

		w.Header().Set("Content-Type", "application-json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)

	}

	product.ID = id

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&product); err != nil {
		errorMessage := map[string]string{"error": "Invalid Payload"}
		body, _ := json.Marshal(errorMessage)

		w.Header().Set("Content-Type", "application-json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
	}

	if err := product.updateProduct(a.DB); err != nil {
		errorMessage := map[string]string{"error": err.Error()}
		body, _ := json.Marshal(errorMessage)

		w.Header().Set("Content-Type", "application-json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
	}

	updatedProduct, _ := json.Marshal(product)
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(updatedProduct)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {}
