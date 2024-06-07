FROM ubuntu:24.04

RUN apt update && \
    apt upgrade -y && \
    apt install -y locales && \
    apt install -y wkhtmltopdf

WORKDIR generator

COPY gohtmltopdf .
COPY .env .

CMD /generator/gohtmltopdf
