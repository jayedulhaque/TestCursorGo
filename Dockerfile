# Build
FROM golang:alpine AS build

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server .

# Run
FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=build /server /server

EXPOSE 8080

ENV HTTP_ADDR=:8080

CMD ["/server"]
