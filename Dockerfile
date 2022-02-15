FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod /app

COPY go.sum /app

RUN go mod download

COPY . /app

RUN go build -o ruller

ENTRYPOINT [ "/app/ruller" ] 

