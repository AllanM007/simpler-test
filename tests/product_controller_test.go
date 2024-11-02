package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/AllanM007/simpler-test/controllers"
	"github.com/AllanM007/simpler-test/models"
	"github.com/AllanM007/simpler-test/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestContainerDB() (*gorm.DB, func(), error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("host=%s port=%s user=user password=password dbname=testdb sslmode=disable", host, port.Port())
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		pgContainer.Terminate(ctx)
	}

	return db, cleanup, nil
}

var db *gorm.DB
var router *gin.Engine
var cleanup func()

// TestMain sets up the test container for all tests and cleans up afterward.
func TestMain(m *testing.M) {
	// Setup: Initialize the PostgreSQL container
	var err error
	db, cleanup, err = SetupTestContainerDB()
	if err != nil {
		log.Fatalf("Could not set up postgres test container: %v", err)
	}

	//initialize gin router
	router = routes.Router(db)

	// Wait for the server to be ready with a ping check
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get("http://localhost:8080/ping")
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// Run tests
	code := m.Run()

	// Teardown the test container
	cleanup()

	// Exit with the code returned by m.Run()
	os.Exit(code)
}

func TestPing(t *testing.T) {

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "pong", recorder.Body.String())
}

var product = controllers.ProductCreateReq{
	Name:        "Test Product",
	Description: "This is a test description for testing product creation",
	Price:       25.50,
	StockLevel:  100,
}

func TestCreateProduct(t *testing.T) {

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

	recorder := httptest.NewRecorder()

	request, err := http.NewRequest(http.MethodGet, "/api/v1/products/1000001", nil)
	if err != nil {
		t.Fatalf("error building request: %v", err)
	}

	router.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
}

func TestUpdateProduct(t *testing.T) {

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

	request, err := http.NewRequest(http.MethodDelete, "/api/v1/products/1", nil)
	if err != nil {
		t.Fatalf("error building request %v", err)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
}
