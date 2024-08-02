# Builder
FROM --platform=linux/amd64 golang:latest AS builder
WORKDIR /app
COPY . .
RUN go work sync && go mod download

# Ad Server Stage
FROM builder AS adserver
WORKDIR /app
RUN go build -o ./adserver/bin ./adserver
EXPOSE 8081
WORKDIR /app/adserver
CMD ["./adserver/bin"]

# Event Server Stage
FROM builder AS eventservice
WORKDIR /app
RUN apt-get update && apt-get install -y librdkafka-dev
RUN go build -o ./eventservice/bin ./eventservice
EXPOSE 8082
WORKDIR /app/eventservice
CMD ["./eventservice/bin"]

# Panel Stage
FROM builder AS panel
WORKDIR /app
RUN go build -o ./panel/bin ./panel
EXPOSE 8085
WORKDIR /app/panel
CMD ["./panel/bin"]

# Publisher Website Stage
FROM builder AS publisherwebsite
WORKDIR /app
RUN go build -o ./publisherwebsite/bin ./publisherwebsite
EXPOSE 8084
WORKDIR /app/publisherwebsite
CMD ["./publisherwebsite/bin"]

FROM builder AS screper
WORKDIR /app
RUN go build -o ./screper/bin ./screper
EXPOSE 8088
WORKDIR /app/screper
CMD ["./screper/bin"]

FROM builder AS reporter
WORKDIR /app
RUN go build -o ./reporter/bin ./reporter
WORKDIR /app/reporter
CMD ["./reporter/bin"]