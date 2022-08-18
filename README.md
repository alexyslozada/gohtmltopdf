# Go HTML to PDF

This is a service that creates a PDF file from a HTML source. The request is made via API REST and the response
is a `[]byte` that contains the PDF file.

## Requisites

We need install `wkhtmltopdf` b/c we are using that library for create PDF files.

```bash
sudo apt install wkhtmltopdf
```

## Installation

We need download de project, configure the `.env` file, compile and run.

Clone the project:

```bash
git clone git@github.com:alexyslozada/gohtmltopdf.git
cd gohtmltopdf
```

Configure the `.env` file:

```bash
cp .env.example .env
# Edit the file with your desire config.
```

Compile and run:
```bash
go mod tidy
go build -o gohtmltopdf cmd/main.go
./gohtmltopdf
```

## Client example

This project has a client example in order to know how to write your own client.
The client is in `client/main.go`. You only need to know how to make a request 
to the server and how to process the response.

If you want to run the client:

```bash
go run client/main.go
```
