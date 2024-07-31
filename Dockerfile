# Builder
FROM golang:latest AS builder
WORKDIR /app

COPY . .
#COPY go.work go.work.sum ./
#
#COPY common ./common
#
#COPY ./adserver ./adserver
#COPY ./eventservice ./eventservice
#COPY ./panel ./panel
#
#COPY ./publisherwebsite ./publisherwebsite
#
#COPY ./adserver/go.mod ./adserver/
#COPY ./adserver/go.sum ./adserver/
#
#COPY ./eventservice/go.mod ./eventservice/
#COPY ./eventservice/go.sum ./eventservice/
#
#COPY ./panel/go.mod ./panel/
#COPY ./panel/go.sum ./panel/
#
#COPY ./publisherwebsite/go.mod ./publisherwebsite/
#COPY ./publisherwebsite/go.sum ./publisherwebsite/

RUN go work sync && go mod download

# Ad Server Stage

FROM builder AS adserver-builder
WORKDIR /app
RUN go build -o adserver /app/adserver

FROM alpine AS adserver
COPY --from=adserver-builder /app/adserver/adserver .
RUN ls > /ls_output.txt
EXPOSE 8081
CMD ["./adserver"]

# Event Server Stage
FROM builder AS eventservice-builder
WORKDIR /app
RUN  go build -o eventservice /app/eventservice

FROM alpine AS eventservice
COPY --from=eventservice-builder /app/eventservice/eventservice /app/eventservice
RUN ls > /ls_output.txt
EXPOSE 8082
CMD ["/app/eventservice"]

# Panel Stage
FROM builder AS panel-builder
WORKDIR /app
RUN go build -o panel /app/panel


FROM alpine AS panel
COPY --from=panel-builder /app/panel/panel .
COPY --from=panel-builder /app/panel/templates ./templates
COPY --from=panel-builder /app/panel/publisher ./publisher
RUN ls > /ls_output.txt
EXPOSE 8085
CMD ["./panel"]

# Publisher Website Stage
FROM builder AS publisherwebsite-builder
WORKDIR /app
RUN go build -o publisherwebsite /app/publisherwebsite

FROM alpine AS publisherwebsite
WORKDIR /app/publisherwebsite
COPY --from=publisherwebsite-builder /app/publisherwebsite/publisherwebsite ./publisherwebsite
COPY --from=publisherwebsite-builder /app/publisherwebsite/html ./html
RUN ls > /ls_output.txt
EXPOSE 8084
CMD ["./publisherwebsite"]