FROM golang:latest

WORKDIR /dnsgoapi

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .

RUN go build -o dnsgoapi cmd/api/main.go

EXPOSE 8080

CMD ["./dnsgoapi", "-port", "8080"]

