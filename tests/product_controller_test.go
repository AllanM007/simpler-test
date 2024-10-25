package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AllanM007/simpler-test/controllers"
	"github.com/AllanM007/simpler-test/models"
	"github.com/AllanM007/simpler-test/routes"
	"github.com/stretchr/testify/assert"
)

var product = models.Product{
	Name:        "Test Product",
	Description: "This is a test description for testing product creation",
	Price:       25.50,
	StockLevel:  100,
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
	router := routes.Router()

	recorder := httptest.NewRecorder()

	jsonValue, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("error marshalling json %v", err)
	}
	request, err := http.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("error buidling request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusCreated, recorder.Code)

}

func TestCreateDuplicateProduct(t *testing.T) {
	router := routes.Router()

	recorder := httptest.NewRecorder()

	jsonValue, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("error marshalling json %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("error buidling request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusConflict, recorder.Code)
}

type ProductsResponse struct {
	Status string                                `json:"status"`
	Data   controllers.ProductsPaginatedResponse `json:"data"`
}

func TestGetProducts(t *testing.T) {
	router := routes.Router()

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/products", nil)
	if err != nil {
		t.Fatalf("error building request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	var products ProductsResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &products)
	if err != nil {
		return
	}

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotEmpty(t, products.Data.Products)
}

type ProductResponse struct {
	Status string                  `json:"status"`
	Data   controllers.ProductData `json:"data"`
}

func TestGetProductById(t *testing.T) {
	router := routes.Router()

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/products/1", nil)
	if err != nil {
		t.Fatalf("error building request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	var product ProductResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &product)
	if err != nil {
		return
	}

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotEmpty(t, product.Data)
}

func TestGetNonExistentProduct(t *testing.T) {
	router := routes.Router()

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/products/1000001", nil)
	if err != nil {
		t.Fatalf("error building request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestUpdateProduct(t *testing.T) {
	router := routes.Router()

	product := models.Product{
		Name:        "Pagani",
		Description: "This is the updated description of koenigsegg to pagani!!",
		Price:       25.40,
	}

	jsonValue, err := json.Marshal(product)
	if err != nil {
		t.Fatalf("error mashalling json %v", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/api/v1/products/1", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("error building request: %v", err)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)

	//test update for non-existent item
	requestNotFound, err := http.NewRequest(http.MethodPut, "/api/v1/products/1000001", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("error building request: %v", err)
	}
	recorderNotFound := httptest.NewRecorder()
	router.ServeHTTP(recorderNotFound, requestNotFound)
	assert.Equal(t, http.StatusNotFound, recorderNotFound.Code)
}

func TestProductSale(t *testing.T) {
	router := routes.Router()

	productSale := controllers.ProductSale{
		Id:    1,
		Count: 200,
	}

	jsonValue, err := json.Marshal(productSale)
	if err != nil {
		t.Fatalf("error mashalling json %v", err)
	}

	request, err := http.NewRequest(http.MethodPut, "/api/v1/products/1/sale", bytes.NewBuffer(jsonValue))
	if err != nil {
		t.Fatalf("error building request %v", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestDeleteProduct(t *testing.T) {
	router := routes.Router()

	request, err := http.NewRequest(http.MethodDelete, "/api/v1/products/1", nil)
	if err != nil {
		t.Fatalf("error building request %v", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
