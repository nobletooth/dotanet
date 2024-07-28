FROM golang:latest AS builder
WORKDIR /app


COPY go.work go.work.sum ./

COPY common ./common

COPY eventservice ./eventservice
COPY panel ./panel
COPY publisherwebsite ./publisherwebsite

COPY ./adserver/go.mod ./adserver/
COPY ./adserver/go.sum ./adserver/

RUN go work sync && go mod download
COPY ./adserver ./adserver
WORKDIR /app/adserver

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o adserver .
FROM scratch
COPY --from=builder /app/adserver/adserver .
EXPOSE 8081
CMD ["./adserver"]
