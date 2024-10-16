package main

import (
	"os"

	"github.com/AllanM007/simpler-test/initializers"
	"github.com/AllanM007/simpler-test/routes"
)

func init() {
	os.Setenv("TZ", "Africa/Nairobi")
	initializers.LoadEnvVariables()
	initializers.InitDb()
}

func main() {
	routes.Router()
}
