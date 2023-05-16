# syntax=docker/dockerfile:1

# Builder image
FROM golang:1.20 AS BUILDER
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY msgme.go .
RUN CGO_ENABLED=0 go build -o msgme .

# App image
FROM scratch
WORKDIR /
COPY .env /
COPY --from=BUILDER /app/msgme /
ENV PORT 3000
EXPOSE 3000
CMD ["/msgme"]
