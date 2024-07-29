# Builder
FROM golang:latest AS builder
WORKDIR /app

COPY go.work go.work.sum ./
COPY common ./common

COPY ./adserver ./adserver
COPY ./eventservice ./eventservice
COPY ./panel ./panel

COPY ./publisherwebsite ./publisherwebsite

COPY ./adserver/go.mod ./adserver/
COPY ./adserver/go.sum ./adserver/

COPY ./eventservice/go.mod ./eventservice/
COPY ./eventservice/go.sum ./eventservice/

COPY ./panel/go.mod ./panel/
COPY ./panel/go.sum ./panel/

COPY ./publisherwebsite/go.mod ./publisherwebsite/
COPY ./publisherwebsite/go.sum ./publisherwebsite/

RUN go work sync && go mod download

# Ad Server Stage

FROM builder AS adserver-builder
WORKDIR /app/adserver
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o adserver .

FROM scratch AS adserver
COPY --from=adserver-builder /app/adserver/adserver .
EXPOSE 8081
CMD ["./adserver"]

# Event Server Stage
FROM builder AS eventservice-builder
WORKDIR /app/eventservice
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o eventservice .

FROM scratch AS eventservice
COPY --from=eventservice-builder /app/eventservice/eventservice .
EXPOSE 8082
CMD ["./eventservice"]

# Panel Stage
FROM builder AS panel-builder
WORKDIR /app/panel
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o panel .


FROM scratch AS panel
COPY --from=panel-builder /app/panel/panel .
COPY --from=panel-builder /app/panel/templates .
COPY --from=panel-builder /app/panel/publisher .
EXPOSE 8085
CMD ["./panel"]

# Publisher Website Stage
FROM builder AS publisherwebsite-builder
WORKDIR /app/publisherwebsite
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o publisherwebsite .

FROM alpine AS publisherwebsite
COPY --from=publisherwebsite-builder /app/publisherwebsite/publisherwebsite .
COPY --from=publisherwebsite-builder /app/publisherwebsite/html .

EXPOSE 8084
CMD ["./publisherwebsite"]