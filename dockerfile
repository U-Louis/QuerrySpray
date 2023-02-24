FROM golang:1.20.1-alpine3.16

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy

ENTRYPOINT ["sh","./build.dev.sh"]