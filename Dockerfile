# syntax=docker/dockerfile:1

FROM golang:1.20
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY msgme.go ./
RUN go build -o /msgme
COPY .env ./
ENV PORT 3000
CMD ["/msgme"]
EXPOSE 3000
