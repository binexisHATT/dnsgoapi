all:
	go build -o dnsapi cmd/dnsgoapi/main.go
clean:
	rm dnsapi
