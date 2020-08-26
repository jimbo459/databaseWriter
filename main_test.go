package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	. "github.com/jimbo459/databaseWriter"
)

var a App

func TestMain(m *testing.M) {
	a.Initialise(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array, got %s", body)
	}
}

func TestGetNonExistentProduct(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		t.Errorf("Error unmarshaling JSON body - %v", response.Body.Bytes())
	}

	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}

}

func TestCreateProduct(t *testing.T) {
	clearTable()

	var jsonString = []byte(`{"name":"test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var responseBody map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	if responseBody["Name"] != "test product" {
		t.Errorf("Expected name to be 'test product', got %s", responseBody["Name"])
	}
	if responseBody["Price"] != 11.22 {
		t.Errorf("Expected price to be '11.22', got %v", responseBody["Price"])
	}
	if responseBody["ID"] != 1.0 {
		t.Errorf("Expected id to be '1', got %v", responseBody["ID"])
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	clearTable()
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &responseBody)

	body := []byte(`{"Name": "Updated Product", "price": 11.22}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(body))
	updatedResponse := executeRequest(req)

	checkResponseCode(t, http.StatusOK, updatedResponse.Code)
	var secondResponse map[string]interface{}
	json.Unmarshal(updatedResponse.Body.Bytes(), &secondResponse)

	if responseBody["Name"] == secondResponse["Name"] {
		t.Errorf("Expected Product name: %v, to equal 'Updated Product'", secondResponse["name"])
	}

	if responseBody["ID"] != secondResponse["ID"] {
		t.Errorf("Expected Product ID %v, to equal 1", secondResponse["ID"])
	}

	if responseBody["Price"] == secondResponse["Price"] {
		t.Errorf("Expected Product Price %v, to equal 11.22", secondResponse["Price"])
	}

}

// func TestDeleteProduct(t *testing.T) {
// 	clearTable()
// 	addProducts(1)

// 	req, _ := http.NewRequest("GET", "/product/1", nil)
// 	response := executeRequest(req)

// 	checkResponseCode(t, http.StatusOK, response.Code)

// 	req, _ = http.NewRequest("DELETE", "/product/1", nil)
// 	response = executeRequest(req)

// 	checkResponseCode(t, http.StatusOK, response.Code)

// 	req, _ = http.NewRequest("GET", "/product/1", nil)
// 	response = executeRequest(req)

// 	checkResponseCode(t, http.StatusNotFound, response.Code)

// }

func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	_, err := a.DB.Exec("DELETE FROM products")
	if err != nil {
		fmt.Printf("Could not clear table")
	}
	_, err = a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	if err != nil {
		fmt.Printf("Could not alter table sequence")
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d to equal %d", actual, expected)
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
	(
		id SERIAL,
		name TEXT NOT NULL,
		price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
		CONSTRAINT products_pkey PRIMARY KEY (id)
	)`
