FROM golang:1.22-alpine AS BUILD

WORKDIR /app

COPY go.mod /app

COPY go.sum /app

RUN go mod download

COPY . /app

RUN go build -o ruller

FROM alpine:3.19

COPY --from=BUILD /app/ruller /bin/

CMD [ "ruller" ] 



