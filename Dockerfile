FROM golang:1.21-alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY src/go.mod src/go.sum config.env ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]