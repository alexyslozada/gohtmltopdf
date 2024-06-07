package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const html = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
    <style>
	.page-break {
		page-break-before: always;
	}
    </style>
</head>
<body>
    <h1 style="font-family: sans-serif; font-size: 4rem;">Hola mundo</h1>
    <p style="font-family: serif; font-size: 0.8rem;">Lorem ipsum dolor sit amet, consectetur adipisicing elit. Consequuntur doloribus nam nemo odio quam quo repudiandae sunt tenetur. Ab aliquam aut beatae dicta minus nesciunt quam quidem similique, temporibus vero.</p>

	<div class="page-break"></div>

    <h1 style="font-family: sans-serif; font-size: 4rem;">Hola mundo</h1>
    <p style="font-family: serif; font-size: 0.8rem;">Lorem ipsum dolor sit amet, consectetur adipisicing elit. Consequuntur doloribus nam nemo odio quam quo repudiandae sunt tenetur. Ab aliquam aut beatae dicta minus nesciunt quam quidem similique, temporibus vero.</p>
</body>
</html>
`

const (
	InternalCodeKey = "INTERNAL_CODE"
	PortKey         = "HTTP_PORT"
)

type Config struct {
	internalCode string
	port         string
}

type request struct {
	Data string `json:"data"`
}

func main() {
	envFilePath := flag.String("envfile", ".env", "Path to the .env file. Default: .env")
	flag.Parse()

	err := loadEnvs(*envFilePath)
	if err != nil {
		log.Fatalf("Couldn´t read de env file in %q path, error: %v", *envFilePath, err)
	}

	src := request{Data: html}
	data, err := json.Marshal(src)
	if err != nil {
		log.Fatalf("error marshaling json: %v", err)
	}

	config := parseEnvToConfig()

	url := fmt.Sprintf("http://localhost:%s/html-to-pdf/%s", config.port, config.internalCode)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		log.Fatalf("error creating the request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error making the request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading the body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d, body: %s", resp.StatusCode, string(body))
	}

	respData := make(map[string][]byte)
	err = json.Unmarshal(body, &respData)
	if err != nil {
		log.Fatalf("error unmarshaling body: %v", err)
	}

	err = writeFile(respData["data"])
	if err != nil {
		log.Fatalf("error writing file: %v", err)
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

// writeFile is only for test proposes. --*Don´t use it*--
func writeFile(data []byte) error {
	return os.WriteFile("test.pdf", data, 0666)
}
