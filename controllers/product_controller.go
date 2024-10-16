package controllers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/AllanM007/simpler-test/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductInput struct {
	Name        string  `json:"name"         binding:"required"`
	Description string  `json:"description"  binding:"required"`
	Price       float64 `json:"price"        binding:"required"`
	// StockLevel  int     `json:"stockLevel"   binding:"required"`
}

type ProductHandler struct {
	DB *gorm.DB
}

func ProductsRepository(db *gorm.DB) *ProductHandler {
	return &ProductHandler{
		DB: db,
	}
}

// @Tags Client-side
// @Summary Products
// @Description
// @Accept  json
// @Produce json
// @Success 200 {object}
// @Router /api/v1/create-product [post]
func (p ProductHandler) CreateProduct(ctx *gin.Context) {

	var product ProductInput
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newProduct := models.Product{
		Name:        product.Name,
		Description: product.Description,
	}

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

// @Tags Client-side
// @Summary Products
// @Description
// @Param page query string false "Number of page"
// @Param limit query string false "count in a page"
// @Accept  json
// @Produce json
// @Success 200 {object}
// @Router /api/v1/products [get]
func (p ProductHandler) GetProducts(ctx *gin.Context) {

	var products []models.Product
	result := p.DB.Find(&products)
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
		}

		data = append(data, product)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "data": data})
}

// @Tags Client-side
// @Summary Products
// @Description
// @Accept  json
// @Produce json
// @Success 200 {object}
// @Router /api/v1/product/{id} [get]
func (p ProductHandler) GetProductById(ctx *gin.Context) {
	productId := ctx.Param("id")

	var product models.Product
	result := p.DB.First(&product).Where("id = ?", productId)
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
		"created_at":  product.CreatedAt,
	}

	data = append(data, productObject)

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "data": data})
}

// @Tags Client-side
// @Summary Products
// @Description
// @Accept  json
// @Produce json
// @Success 200 {object} PageableLibraryResponse
// @Router /api/v1/update-product [put]
func (p ProductHandler) UpdateProduct(ctx *gin.Context) {

	productId := ctx.Param("id")

	var updateProduct ProductInput
	if err := ctx.ShouldBindJSON(&updateProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "BAD_REQUEST", "error": err.Error()})
		return
	}

	log.Println(productId)

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
		if strings.Contains(result.Error.Error(), "duplicate key value") {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{"status": "DUPLICATE_ENTITY", "error": "Duplicate conflict while updating product!"})
			return
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Product updated succesfully!"})
}

// @Tags Client-side
// @Summary Products
// @Description
// @Accept  json
// @Produce json
// @Success 200 {object}
// @Router /api/v1/delete-product/{id} [delete]
func (p *ProductHandler) DeleteProduct(ctx *gin.Context) {
	productId := ctx.Param("id")
	log.Println(productId)

	var product models.Product
	result := p.DB.Where("id = ?", productId).Delete(&product)
	log.Println(result)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": "NOT_FOUND", "message": "Product not found!!"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Product deleted successfully!"})
}
