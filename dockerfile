FROM golang:1.20.1-alpine3.16

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

CMD ["go", "run", "main.go"]

EXPOSE 8085