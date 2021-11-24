FROM golang:1.14.6-alpine3.12 as builder
COPY go.mod go.sum /go/src/github.com/GlobantObrikosina/golang-rest-api/
WORKDIR /go/src/github.com/GlobantObrikosina/golang-rest-api
RUN go mod download
COPY . /go/src/github.com/GlobantObrikosina/golang-rest-api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/golang-rest-api github.com/GlobantObrikosina/golang-rest-api

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/GlobantObrikosina/golang-rest-api/build/golang-rest-api /usr/bin/golang-rest-api
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/golang-rest-api"]