package main

import (
	"os"

	_ "github.com/AllanM007/simpler-test/docs"
	"github.com/AllanM007/simpler-test/initializers"
	"github.com/AllanM007/simpler-test/routes"
)

func init() {
	os.Setenv("TZ", "Africa/Nairobi")
	initializers.LoadEnvVariables()
}

func main() {
	routes.Router()
}
