package controllers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/AllanM007/simpler-test/models"
	"github.com/AllanM007/simpler-test/utilities"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductCreateInput struct {
	Name        string  `json:"name"         binding:"required"`
	Description string  `json:"description"  binding:"required"`
	Price       float64 `json:"price"        binding:"required"`
	StockLevel  int     `json:"stockLevel"   binding:"required"`
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
// @Summary Create products
// @Tags products
// @Description create new product
// @Accept  json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/v1/create-product [post]
func (p ProductHandler) CreateProduct(ctx *gin.Context) {
	log.Println("i was hit")

	var product ProductCreateInput
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("i was hit again")

	newProduct := models.Product{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		StockLevel:  product.StockLevel,
	}

	log.Println("i was hit again again")

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

	log.Println("i was hit again again again")

	ctx.JSON(http.StatusCreated, gin.H{"status": "OK", "message": "Product created successfully!"})

}

// List Products
// @Summary Get products
// @Description get all products
// @Tags products
// @Param page query string false "Number of page"
// @Param limit query string false "count in a page"
// @Accept  json
// @Produce json
// @Success 200
// @Failure 404
// @Failure 500
// @Router /api/v1/products [get]
func (p ProductHandler) GetProducts(ctx *gin.Context) {
	page, limit, err := utilities.GetPagingData(ctx)
	if err != nil {
		return
	}

	var products []models.Product
	result := p.DB.Limit(limit).Offset(utilities.GetOffset(page, limit)).Order("id DESC").Find(&products)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "NOT_FOUND", "message": "Products not found!!"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": result.Error})
		return
	}

	var data = []map[string]interface{}{}

	for i := 0; i < len(products); i++ {
		product := map[string]interface{}{
			"id":          products[i].ID,
			"name":        products[i].Name,
			"description": products[i].Description,
			"stock":       products[i].StockLevel,
		}

		data = append(data, product)
	}
	log.Println(data)

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "data": data})
}

// GetProductById godoc
// @Summary Get Product By Id
// @Description get product by id
// @Tags products
// @Accept  json
// @Produce json
// @Success 200
// @Failure 404
// @Failure 500
// @Router /api/v1/product/{id} [get]
func (p ProductHandler) GetProductById(ctx *gin.Context) {
	productId := ctx.Param("id")

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

	var data = []map[string]interface{}{}

	productObject := map[string]interface{}{
		"id":          product.ID,
		"name":        product.Name,
		"description": product.Description,
		"stock":       product.StockLevel,
		"created_at":  product.CreatedAt,
	}

	data = append(data, productObject)
	log.Println(data)

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "data": data})
}

type ProductUpdateInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	StockLevel  int     `json:"stockLevel"`
}

// UpdateProduct godoc
// @Summary Update A Product
// @Description update a product by id
// @Tags products
// @Accept  json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /api/v1/update-product/{id} [put]
func (p ProductHandler) UpdateProduct(ctx *gin.Context) {

	productId := ctx.Param("id")

	var updateProduct ProductUpdateInput
	if err := ctx.ShouldBindJSON(&updateProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "BAD_REQUEST", "error": err.Error()})
		return
	}

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
// @Summary Update A Product
// @Description update a product by id
// @Tags products
// @Accept  json
// @Produce json
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
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

	if productSale.Count > product.StockLevel {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "FORBIDDEN", "message": "Stock level lower than purchase quantity"})
		return
	}

	product.StockLevel = product.StockLevel - productSale.Count

	updateResult := p.DB.Save(&product)
	if updateResult.Error != nil {

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Product sale successful!"})

}

// DeleteProduct godoc
// @Summary Delete A Product
// @Description delete product by id
// @Tags products
// @Accept  json
// @Produce json
// @Success 200
// @Failure 404
// @Failure 500
// @Router /api/v1/delete-product/{id} [delete]
func (p *ProductHandler) DeleteProduct(ctx *gin.Context) {
	productId := ctx.Param("id")
	log.Println(productId)

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
