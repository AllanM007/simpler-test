package controllers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AllanM007/simpler-test/helpers"
	"github.com/AllanM007/simpler-test/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type InternalErrorResponse struct {
	Error string `json:"error"`
}

type ProductCreateReq struct {
	Name        string  `json:"name"         binding:"required"`
	Description string  `json:"description"  binding:"required"`
	Price       float64 `json:"price"        binding:"required"`
	StockLevel  int     `json:"stock"        `
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
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/create-product [post]
func (p ProductHandler) CreateProduct(ctx *gin.Context) {

	var product ProductCreateReq
	if err := ctx.ShouldBindJSON(&product); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
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
	CurrentPage int    `json:"current_page"`
	Limit       int    `json:"limit"`
	Total       *int64 `json:"total_products"`
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": result.Error})
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
	var count *int64
	_ = p.DB.Find(&products).Count(count)

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
// @Router /api/v1/product/{id} [get]
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
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": result.Error})
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
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/update-product/{id} [put]
func (p ProductHandler) UpdateProduct(ctx *gin.Context) {

	productId := ctx.Param("id")

	var updateProduct ProductUpdateReq
	if err := ctx.ShouldBindJSON(&updateProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "BAD_REQUEST", "error": err.Error()})
		return
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

	product.Name = updateProduct.Name
	product.Description = updateProduct.Description

	//update existing products
	updateResult := p.DB.Save(&product)
	if updateResult.Error != nil {
		if strings.Contains(updateResult.Error.Error(), "duplicate key value") {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"status": "DUPLICATE_ENTITY", "error": "Duplicate conflict while updating product!"})
			return
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Product updated succesfully!"})
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
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/product-sale [put]
func (p *ProductHandler) ProductSale(ctx *gin.Context) {

	var productSale ProductSale
	if err := ctx.ShouldBindJSON(&productSale); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "BAD_REQUEST", "error": err.Error()})
		return
	}

	var product models.Product
	result := p.DB.Where("id = ?", productSale.Id).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "NOT_FOUND", "message": "Product not found!!"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	//check if product stock is less than sale quantity
	if productSale.Count > product.StockLevel {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "FORBIDDEN", "message": "Stock level lower than purchase quantity"})
		return
	}

	// dedcut sale quantity from product stock
	product.StockLevel = product.StockLevel - productSale.Count

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
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} InternalErrorResponse
// @Router /api/v1/delete-product/{id} [delete]
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
