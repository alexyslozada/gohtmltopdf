package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/alexyslozada/gohtmltopdf"
)

func main() {
	e := echo.New()
	e.POST("/html-to-pdf", gohtmltopdf.NewHandler().CreatePDF)
	err := e.Start(":9632")
	if err != nil {
		log.Fatalf("Couldn´t start the server, err: %v", err)
	}
}
