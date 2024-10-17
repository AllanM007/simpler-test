package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/AllanM007/simpler-test/controllers"
	"github.com/AllanM007/simpler-test/initializers"
	"github.com/AllanM007/simpler-test/models"
	"github.com/AllanM007/simpler-test/routes"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// initialize database connection for tests
func Setup() {
	err := godotenv.Load(filepath.Join("../", ".env"))
	if err != nil {
		log.Println("error loading env file")
		log.Fatal(err)
	}
	initializers.InitDb()
}

func TestPing(t *testing.T) {
	router := routes.Router()

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "pong", recorder.Body.String())
}

func TestCreateProduct(t *testing.T) {

	Setup()
	router := routes.Router()

	recorder := httptest.NewRecorder()

	productId := time.Now().UnixNano()
	product := models.Product{
		ID:          uint64(productId),
		Name:        "Test Product",
		Description: "This is a test description for testing product creation",
		Price:       25.50,
		StockLevel:  100,
	}
	jsonValue, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("error mashalling json %v", err)
	}
	request, err := http.NewRequest(http.MethodPost, "/api/v1/create-product", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("error buidling request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestGetProducts(t *testing.T) {

	Setup()
	router := routes.Router()

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/products", nil)
	if err != nil {
		t.Fatalf("error building request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	var products []models.Product
	err = json.Unmarshal(recorder.Body.Bytes(), &products)
	if err != nil {
		return
	}

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotEmpty(t, products)
}

func TestGetProductById(t *testing.T) {

	Setup()
	router := routes.Router()

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/product/8", nil)
	if err != nil {
		t.Fatalf("error building request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	var product models.Product
	err = json.Unmarshal(recorder.Body.Bytes(), &product)
	if err != nil {
		return
	}

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotEmpty(t, product)
}

func TestGetNonExistentProduct(t *testing.T) {
	Setup()
	router := routes.Router()

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/product/1000001", nil)
	if err != nil {
		t.Fatalf("error building request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	var product models.Product
	err = json.Unmarshal(recorder.Body.Bytes(), &product)
	if err != nil {
		return
	}

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Empty(t, product)
}

func TestUpdateProduct(t *testing.T) {
	Setup()

	router := routes.Router()

	product := models.Product{
		ID:          8,
		Name:        "Pagani",
		Description: "This is the updated description of koenigsegg to pagani!!",
		Price:       25.40,
	}

	jsonValue, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("error mashalling json %v", err)
	}

	request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/update-product/%d", product.ID), bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("error building request: %v", err)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)

	//test update for non-existent item
	requestNotFound, err := http.NewRequest(http.MethodPut, "/api/v1/update-product/1000001", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("error building request: %v", err)
	}
	recorderNotFound := httptest.NewRecorder()
	router.ServeHTTP(recorderNotFound, requestNotFound)
	assert.Equal(t, http.StatusNotFound, recorderNotFound.Code)
}

func TestDeleteProduct(t *testing.T) {

	Setup()
	router := routes.Router()

	request, err := http.NewRequest(http.MethodDelete, "/api/v1/delete-product/8", nil)
	if err != nil {
		t.Fatalf("error vuilding request %v", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestProductSale(t *testing.T) {
	Setup()

	router := routes.Router()

	productSale := controllers.ProductSale{
		Id:    8,
		Count: 200,
	}

	jsonValue, err := json.Marshal(productSale)
	if err != nil {
		t.Fatalf("error mashalling json %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/api/v1/product-sale", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("error vuilding request %v", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}
