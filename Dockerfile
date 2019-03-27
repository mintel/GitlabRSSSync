FROM golang:latest AS builder
RUN mkdir /app
COPY go.mod /app/
WORKDIR /app
RUN go mod download
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rss_sync .

FROM scratch
COPY --from=builder /app/rss_sync /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app
CMD ["/app/rss_sync"]