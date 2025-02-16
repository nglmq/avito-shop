FROM golang:1.23

WORKDIR /app
COPY . .

RUN go build -o /build ./cmd/shop \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]