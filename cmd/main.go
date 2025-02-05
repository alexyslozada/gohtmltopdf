package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/alexyslozada/gohtmltopdf"
)

const (
	InternalCodeKey = "INTERNAL_CODE"
	PortKey         = "HTTP_PORT"
)

type Config struct {
	internalCode string
	port         string
}

func main() {
	envFilePath := flag.String("envfile", ".env", "Path to the .env file. Default: .env")
	flag.Parse()

	err := loadEnvs(*envFilePath)
	if err != nil {
		log.Fatalf("Couldn´t read de env file in %q path, error: %v", *envFilePath, err)
	}

	config := parseEnvToConfig()

	e := echo.New()
	gohtmltopdf.Router(e, config.internalCode)

	err = e.Start(fmt.Sprintf(":%s", config.port))
	if err != nil {
		log.Fatalf("Couldn´t start the server, err: %v", err)
	}
}

func loadEnvs(envFilePath string) error {
	return godotenv.Load(envFilePath)
}

func parseEnvToConfig() Config {
	internalCode := os.Getenv(InternalCodeKey)
	port := os.Getenv(PortKey)

	return Config{
		internalCode: internalCode,
		port:         port,
	}
}
