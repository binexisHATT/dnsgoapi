FROM golang:latest

WORKDIR /go/src/dnsgoapi
COPY . .

RUN go get -u github.com/miekg/dns
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/russross/blackfriday

RUN go install cmd/dnsgoapi 

CMD ["dnsgoapi"]

