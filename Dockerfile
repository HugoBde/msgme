# syntax=docker/dockerfile:1

# Builder image
FROM golang:1.20-alpine AS BUILDER
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY msgme.go .
RUN CGO_ENABLED=0 go build -o msgme .
RUN apk add --no-cache ca-certificates

# App image
FROM scratch
WORKDIR /
COPY .env .
COPY --from=BUILDER /app/msgme .
WORKDIR /etc/ssl/certs
COPY --from=BUILDER /etc/ssl/certs/ca-certificates.crt .
WORKDIR /
ENV PORT 3000
EXPOSE 3000
CMD ["/msgme"]
