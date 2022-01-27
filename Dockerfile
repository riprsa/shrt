FROM golang:alpine as builder

RUN go version
RUN apk add git

COPY ./ /github.com/hararudoka/shorter
WORKDIR /github.com/hararudoka/shorter

RUN go mod download && go get -u ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./main.go

#lightweight docker container with binary
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=0 /github.com/hararudoka/shorter/.bin/app .

EXPOSE 80/tcp
EXPOSE 443/tcp

CMD [ "./app"]