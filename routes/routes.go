package routes

import (
	"net/http"

	"github.com/AllanM007/simpler-test/controllers"
	"github.com/AllanM007/simpler-test/initializers"
	"github.com/AllanM007/simpler-test/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	app := gin.Default()

	app.GET("/ping", func(ctx *gin.Context) {
		ctx.String(
			http.StatusOK,
			"pong",
		)
	})

	// enable cors middleware to apply to all routes
	app.Use(middleware.CORSMiddleware())

	db := initializers.InitDb()
	ProductsRepo := controllers.ProductsRepository(db)

	app.POST("/api/v1/create-product", ProductsRepo.CreateProduct)
	app.GET("/api/v1/products", ProductsRepo.GetProducts)
	app.GET("/api/v1/product/:id", ProductsRepo.GetProductById)
	app.PUT("/api/v1/update-product/:id", ProductsRepo.UpdateProduct)
	app.DELETE("/api/v1/delete-product/:id", ProductsRepo.DeleteProduct)

	app.POST("/api/v1/product-sale")

	app.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	app.Run(":8080")

	return app
}
