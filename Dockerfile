FROM golang:alpine as builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin ./cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /build .

EXPOSE 80/tcp 443/tcp

ENTRYPOINT ["/app/bin"]