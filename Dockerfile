FROM golang:1.22
WORKDIR /app
COPY . .
RUN go mod download
EXPOSE 8080
ENTRYPOINT exec go run cmd/main.go


