FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o calculator-api


WORKDIR /app

COPY --from=builder /app/calculator-api .

EXPOSE 8080

CMD ["./calculator-api"]
