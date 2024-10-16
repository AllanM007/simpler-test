package routes

import (
	"net/http"

	"github.com/AllanM007/simpler-test/controllers"
	"github.com/AllanM007/simpler-test/initializers"
	"github.com/AllanM007/simpler-test/middleware"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	app := gin.Default()

	// enable cors middleware to apply to all routes
	app.Use(middleware.CORSMiddleware())

	app.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"message": "pong",
			})
	})

	db := initializers.InitDb()
	ProductsRepo := controllers.ProductsRepository(db)

	app.POST("/api/v1/create-product", ProductsRepo.CreateProduct)
	app.GET("/api/v1/products", ProductsRepo.GetProducts)
	app.GET("/api/v1/product/:id", ProductsRepo.GetProductById)
	app.PUT("/api/v1/update-product/:id", ProductsRepo.UpdateProduct)
	app.DELETE("/api/v1/delete-product/:id", ProductsRepo.DeleteProduct)

	app.Run(":8080")

	return app
}
