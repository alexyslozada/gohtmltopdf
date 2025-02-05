package gohtmltopdf

import "github.com/labstack/echo/v4"

func Router(e *echo.Echo, internalCode string) {
	handler := NewHandler()
	e.GET("/health", handler.Health)
	e.POST("/html-to-pdf", handler.ValidateInternalCode(handler.CreateHTMLToPDF, internalCode))
	e.POST("/dian-form-220", handler.ValidateInternalCode(handler.CreateDianForm220, internalCode))
}
