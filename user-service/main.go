package main

import (
	"log"
	"user-service/config"
	"user-service/router"

	"github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()

    config.InitDB()

    router.InitRoutes(e)

    if err := e.Start(":5002"); err != nil {
        log.Fatal(err)
    }
}
