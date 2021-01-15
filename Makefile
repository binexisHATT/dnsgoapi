all:
	go build -o dnsgoapi cmd/api/main.go
clean:
	rm dnsapi
