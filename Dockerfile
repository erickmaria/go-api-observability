FROM golang:1.23.5 AS builder

ARG MODE

RUN CGO_ENABLED=1
WORKDIR /app
COPY go.mod go.sum config.yaml ./
RUN go mod download
COPY . .

RUN go build -o random-number-api cmd/${MODE}/main.go


# FROM scratch
# COPY --from=builder /app/go-api /app/go-api
# COPY --from=builder /app/config.yaml /app/config.yaml
ENTRYPOINT ["sh", "-c", "/app/random-number-api"]