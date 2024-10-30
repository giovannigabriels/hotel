package main

import (
	"api-gateway/router"
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	router.InitRoutes(e)

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting API Gateway on port %s\n", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
