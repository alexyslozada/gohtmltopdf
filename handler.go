package gohtmltopdf

import (
	"bytes"
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

// request is the structure that the client has to send to parse to pdf
type request struct {
	// TODO agregar propiedades como: tamaño de página, número de páginas, etc

	// Data must be a string with html format.
	Data string `json:"data"`
}

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) CreatePDF(c echo.Context) error {
	req := request{}
	err := c.Bind(&req)
	if err != nil {
		errMsg := map[string]string{"msg": "can't bind request", "error": err.Error()}
		c.Logger().Error(errMsg)
		return c.JSON(http.StatusBadRequest, errMsg)
	}

	src := bytes.NewBufferString(req.Data)
	gen := NewGenerator(src)
	pdf, err := gen.run(context.Background())
	if err != nil {
		errMsg := map[string]string{"msg": "can't create the PDF", "error": err.Error()}
		c.Logger().Error(errMsg)
		return c.JSON(http.StatusInternalServerError, errMsg)
	}

	return c.JSON(http.StatusOK, map[string][]byte{"data": pdf})
}
