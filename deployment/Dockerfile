FROM golang:1.18-alpine3.15 as builder

WORKDIR /go/api
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/main.go

FROM scratch
COPY --from=builder /go/api/main /go/bin/main
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/go/bin/main"]
