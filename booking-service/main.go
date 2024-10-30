package main

import (
	"booking-service/config"
	"booking-service/router"
	"log"

	"github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()

    config.InitDB()

    router.InitRoutes(e)

    if err := e.Start(":5001"); err != nil {
        log.Fatal(err)
    }
}
