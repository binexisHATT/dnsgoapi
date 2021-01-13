all:
	go build -o dnsapi cmd/api/main.go
clean:
	rm dnsapi
