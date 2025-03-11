# The build stage
FROM golang:1.24 as builder
WORKDIR /api
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/*.go

# The run stage
FROM scratch
WORKDIR /api
# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/api .
EXPOSE 8080
CMD ["./api"]