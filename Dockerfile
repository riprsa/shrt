FROM golang:alpine as builder

WORKDIR /build

COPY go.mod go.sum docker-compose.yaml ./

RUN go mod download

COPY . .

RUN go build

FROM alpine

WORKDIR /app

COPY --from=builder /build/shorter .

EXPOSE 80/tcp 443/tcp

ENTRYPOINT ["/app/shorter"]