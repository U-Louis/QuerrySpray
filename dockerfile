FROM golang:latest

WORKDIR /app

RUN go install github.com/gin-gonic/gin@latest

VOLUME ["/app"]

RUN go build -o QuerrySpray

EXPOSE 8085

CMD ["./QuerrySpray"]
