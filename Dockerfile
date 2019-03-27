FROM golang:latest AS builder
RUN mkdir /app
COPY go.mod /app/
WORKDIR /app
RUN go mod download
COPY . /app
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o rss_sync .
WORKDIR /app
CMD ["/app/rss_sync"]