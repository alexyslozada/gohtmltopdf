FROM --platform=linux/amd64 ubuntu:24.04

RUN apt update && \
  apt upgrade -y && \
  apt install -y locales && \
  apt install -y wkhtmltopdf

WORKDIR /genpdf

COPY gohtmltopdf .
COPY .env .

CMD ["/genpdf/gohtmltopdf"]
