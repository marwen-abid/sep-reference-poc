FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o /sep-reference ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /sep-reference /usr/local/bin/sep-reference
EXPOSE 8080
CMD ["/usr/local/bin/sep-reference"]
