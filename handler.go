package gohtmltopdf

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) CreateHTMLToPDF(c echo.Context) error {
	req := requestHTML{}
	err := c.Bind(&req)
	if err != nil {
		errMsg := map[string]string{"msg": "can't bind requestHTML", "error": err.Error()}
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

func (h Handler) CreateDianForm220(c echo.Context) error {
	req := requestDIANForm220{}
	err := c.Bind(&req)
	if err != nil {
		errMsg := map[string]string{"msg": "can't bind requestDIANForm220", "error": err.Error()}
		c.Logger().Error(errMsg)
		return c.JSON(http.StatusBadRequest, errMsg)
	}

	// If we need to debug the performance, we can set the query param debug=true
	isDebug := false
	isDebugStr := c.QueryParam("debug")
	if strings.EqualFold(isDebugStr, "true") {
		isDebug = true
	}

	dian := NewDIAN(isDebug)
	pdf, err := dian.CreateDIANForm220(req.Data)
	if err != nil {
		if errors.As(err, &ErrorProcess{}) {
			errMsg := map[string]string{"msg": "can't create the PDF", "error": err.Error()}
			c.Logger().Error(errMsg)
			return c.JSON(http.StatusBadRequest, errMsg)
		}

		errMsg := map[string]string{"msg": "can't create the PDF", "error": err.Error()}
		c.Logger().Error(errMsg)
		return c.JSON(http.StatusInternalServerError, errMsg)
	}

	return c.JSON(http.StatusOK, map[string][]byte{"data": pdf})
}

func (h Handler) Health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"date": time.Now().String()})
}

const ParamInternalCode = "x-internalcode"

// ValidateInternalCode to validate the internal code
func (h Handler) ValidateInternalCode(next echo.HandlerFunc, internalCode string) echo.HandlerFunc {
	return func(c echo.Context) error {
		internalReceived := c.Request().Header.Get(ParamInternalCode)
		if internalReceived != internalCode {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "The header x-internal code sent is not valid"})
		}

		return next(c)
	}
}
