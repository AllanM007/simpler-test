package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/AllanM007/simpler-test/helpers"
	"github.com/AllanM007/simpler-test/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type InvalidRequestResponse struct {
	Errors map[string]string `json:"errors"`
	Status string            `json:"status"`
}

type InternalErrorResponse struct {
	Error string `json:"error"`
}

type ProductCreateReq struct {
	Name        string  `json:"name"         binding:"required"`
	Description string  `json:"description"  binding:"required"`
	Price       float64 `json:"price"        binding:"required,gt=0"`
	StockLevel  int     `json:"stock"        binding:"required,gt=0"`
}

type ProductHandler struct {
	DB *gorm.DB
}

func ProductsRepository(db *gorm.DB) *ProductHandler {
	return &ProductHandler{
		DB: db,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Tags products
// @Description create product
// @Accept  json
// @Produce json
// @Param params body ProductCreateReq true "Request's body"
// @Success 200 {object} Response
// @Failure 400 {object} InvalidRequestResponse
// @Failure 404 {object} Response
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/products [post]
func (p ProductHandler) CreateProduct(ctx *gin.Context) {

	var product ProductCreateReq
	if err := ctx.ShouldBindJSON(&product); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := formatValidationError(validationErrors)
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "BAD_REQUEST", "errors": errors})
			return
		}
	}

	newProduct := models.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		StockLevel:  product.StockLevel,
	}

	//insert new product item to database
	result := p.DB.Create(&newProduct)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value") {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"status": "DUPLICATE_ENTITY", "error": "Duplicate conflict while creating product!"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return

	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "OK", "message": "Product created successfully!"})

}

type ProductData struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}

type RequestMeta struct {
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
	Total       int64 `json:"total_products"`
}

type ProductsPaginatedResponse struct {
	Products []ProductData `json:"products"`
	Meta     RequestMeta   `json:"meta"`
}

// Get Products
// @Summary Get products with paging
// @Description get all products
// @Tags products
// @Param page       query string false "Number of page"        default(1)
// @Param limit      query string false "Books count in a page" default(10)
// @Accept  json
// @Produce json
// @Success 200 {object} ProductsPaginatedResponse
// @Failure 404 {object} Response
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/products [get]
func (p ProductHandler) GetProducts(ctx *gin.Context) {
	page, limit, err := helpers.GetPagingData(ctx)
	if err != nil {
		return
	}

	//get products based on pagination parameters
	var products []models.Product
	result := p.DB.Limit(limit).Offset(helpers.GetOffset(page, limit)).Order("id DESC").Find(&products)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "NOT_FOUND", "message": "Products not found!!"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	var data []ProductData

	for i := 0; i < len(products); i++ {
		product := ProductData{
			Id:          products[i].ID,
			Name:        products[i].Name,
			Description: products[i].Description,
			Stock:       products[i].StockLevel,
			CreatedAt:   products[i].CreatedAt,
		}

		data = append(data, product)
	}

	//get all products count
	var count int64
	tx := p.DB.Find(&products).Count(&count)
	if tx.Error != nil {
		return
	}

	meta := RequestMeta{
		CurrentPage: page,
		Limit:       limit,
		Total:       count,
	}

	response := ProductsPaginatedResponse{
		Products: data,
		Meta:     meta,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "data": response})
}

// GetProductById godoc
// @Summary Get product
// @Description get product by id
// @Tags products
// @Param id path int true "Product Id"
// @Accept  json
// @Produce json
// @Success 200 {object} ProductData
// @Failure 404 {object} Response
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/products/{id} [get]
func (p ProductHandler) GetProductById(ctx *gin.Context) {
	productId := ctx.Param("id")

	//get product using id
	var product models.Product
	result := p.DB.Where("id = ?", productId).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "NOT_FOUND", "message": "Product not found!!"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	var data []ProductData

	productObject := ProductData{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Stock:       product.StockLevel,
		CreatedAt:   product.CreatedAt,
	}

	data = append(data, productObject)

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "data": data})
}

type ProductUpdateReq struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	StockLevel  int     `json:"stockLevel"`
}

// UpdateProduct godoc
// @Summary Update product
// @Description update a product by id
// @Tags products
// @Param id path int true "Product Id"
// @Accept  json
// @Produce json
// @Param params body ProductUpdateReq true "Request's body"
// @Success 200 {object} Response
// @Failure 400 {object} InvalidRequestResponse
// @Failure 404 {object} Response
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/products/{id} [put]
func (p ProductHandler) UpdateProduct(ctx *gin.Context) {

	productId := ctx.Param("id")

	var updateProductReq ProductUpdateReq
	if err := ctx.ShouldBindJSON(&updateProductReq); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := formatValidationError(validationErrors)
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "BAD_REQUEST", "errors": errors})
			return
		}
	}

	//get product by id
	var product models.Product
	result := p.DB.Where("id = ?", productId).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "NOT_FOUND", "message": "Product not found!!"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	product.Name = updateProductReq.Name
	product.Description = updateProductReq.Description

	//update existing products
	updateResult := p.DB.Save(&product)
	if updateResult.Error != nil {
		if strings.Contains(updateResult.Error.Error(), "duplicate key value") {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"status": "DUPLICATE_ENTITY", "error": "Duplicate conflict while updating product!"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Product updated successfully!"})
}

type ProductSale struct {
	Id    int `json:"id" binding:"required"`
	Count int `json:"count" binding:"required"`
}

// UpdateProduct godoc
// @Summary Product sale
// @Description  product sale
// @Tags products
// @Accept  json
// @Produce json
// @Param params body ProductSale true "Request's body"
// @Success 200 {object} Response
// @Failure 400 {object} InvalidRequestResponse
// @Failure 404 {object} Response
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/products/{id}/sale [put]
func (p *ProductHandler) ProductSale(ctx *gin.Context) {

	var productSaleReq ProductSale
	if err := ctx.ShouldBindJSON(&productSaleReq); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := formatValidationError(validationErrors)
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "BAD_REQUEST", "errors": errors})
			return
		}
	}

	var product models.Product
	result := p.DB.Where("id = ?", productSaleReq.Id).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "NOT_FOUND", "message": "Product not found!!"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	//check if product stock is less than sale quantity
	if productSaleReq.Count > product.StockLevel {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "FORBIDDEN", "message": "Stock level lower than purchase quantity"})
		return
	}

	// dedcut sale quantity from product stock
	product.StockLevel = product.StockLevel - productSaleReq.Count

	updateResult := p.DB.Save(&product)
	if updateResult.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Product sale successful!"})

}

// DeleteProduct godoc
// @Summary Delete product
// @Description delete product by id
// @Tags products
// @Param id path int true "Product Id"
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} InvalidRequestResponse
// @Failure 404 {object} Response
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/products/{id} [delete]
func (p *ProductHandler) DeleteProduct(ctx *gin.Context) {
	productId := ctx.Param("id")

	//delete product with specified id
	var product models.Product
	result := p.DB.Where("id = ?", productId).Delete(&product)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	} else if result.RowsAffected == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "NOT_FOUND", "message": "Product not found!!"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Product deleted successfully!"})
}

func formatValidationError(errs validator.ValidationErrors) map[string]string {
	errorMessages := make(map[string]string)
	for _, err := range errs {
		field := err.Field()
		switch err.Tag() {
		case "required":
			errorMessages[field] = field + " is required"
		case "min":
			errorMessages[field] = field + " value is too low"
		case "max":
			errorMessages[field] = field + " value is too high"
		default:
			errorMessages[field] = "Invalid value for " + field
		}
	}
	return errorMessages
}
