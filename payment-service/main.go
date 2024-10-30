package main

import (
	"log"
	"payment-service/config"
	"payment-service/router"

	"github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()

    config.InitDB()

    router.InitRoutes(e)

    if err := e.Start(":5003"); err != nil {
        log.Fatal(err)
    }
}
