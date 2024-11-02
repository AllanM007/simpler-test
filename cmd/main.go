package main

import (
	"log"
	"os"

	_ "github.com/AllanM007/simpler-test/docs"
	"github.com/AllanM007/simpler-test/initializers"
	"github.com/AllanM007/simpler-test/routes"
)

// @title           Simpler Test API
// @version         1.0
// @description     This is a product resource microservice RESTful API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      http://localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	os.Setenv("TZ", "Africa/Nairobi")
	initializers.LoadEnvVariables()

	db, err := initializers.ConnectDB()
	if err != nil {
		log.Fatalf("failed to initalize database: %v", err)
	}
	routes.Router(db).Run(":8080")
}
