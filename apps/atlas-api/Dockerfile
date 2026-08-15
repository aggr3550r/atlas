FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o atlas-api .

FROM alpine:3.22

WORKDIR /app

RUN addgroup -S atlas && adduser -S atlas -G atlas

COPY --from=builder /app/atlas-api ./atlas-api

RUN chown atlas:atlas ./atlas-api

USER atlas

EXPOSE 8080

CMD ["./atlas-api"]