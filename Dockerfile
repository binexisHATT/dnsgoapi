FROM golang:latest

WORKDIR /go/src/dnsgoapi
COPY . .

RUN go get -u github.com/miekg/dns
RUN go get -u github.com/gorilla/dns

RUN go install cmd/dnsgoapi 

CMD ["dnsgoapi"]

