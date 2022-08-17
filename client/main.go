package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
</head>
<body>
    <h1 style="font-family: sans-serif; font-size: 4rem;">Hola mundo</h1>
    <p style="font-family: serif; font-size: 0.8rem;">Lorem ipsum dolor sit amet, consectetur adipisicing elit. Consequuntur doloribus nam nemo odio quam quo repudiandae sunt tenetur. Ab aliquam aut beatae dicta minus nesciunt quam quidem similique, temporibus vero.</p>
</body>
</html>
`

type request struct {
	Data string `json:"data"`
}

func main() {
	src := request{Data: html}
	data, err := json.Marshal(src)
	if err != nil {
		log.Fatalf("error marshaling json: %v", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "http://localhost:9632/html-to-pdf", bytes.NewReader(data))
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

	body, err := ioutil.ReadAll(resp.Body)
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

// writeFile is only for test proposes. --*Don´t use it*--
func writeFile(data []byte) error {
	return ioutil.WriteFile("test.pdf", data, 0666)
}
