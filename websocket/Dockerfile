FROM golang:1.21.1 AS binarybuilder
WORKDIR /app
COPY . /app
RUN go mod tidy && \
    go build

FROM debian:bookworm-slim
EXPOSE 8080
RUN apt update && \
    apt install -y ca-certificates
WORKDIR /app
COPY --from=binarybuilder /app/websocket .
CMD ["/app/websocket"]