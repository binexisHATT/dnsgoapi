package main

import (
	"fmt"
	"log"
	"net/http"
    "flag"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"

	"dnsgoapi/dnsgoapi"
)

func main() {
    port := flag.Int("port", 8080, "The port for the server to listen on")
    flag.Parse()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        input := []byte(`
	 # dnsgoapi
	 ### An API for making simple DNS queries with various free public DNS server support
	 `
	)
        fmt.Fprintf(w, string(blackfriday.MarkdownCommon(input)))
	}).Methods("GET")

	r.HandleFunc("/a/{PublicDNS}/{q}", dnsgoapi.QueryA).Methods("GET")
	r.HandleFunc("/aaaa/{PublicDNS}/{q}", dnsgoapi.QueryQuadA).Methods("GET")
	r.HandleFunc("/cname/{PublicDNS}/{q}", dnsgoapi.QueryCNAME).Methods("GET")

    l := fmt.Sprintf(":%d", *port)
    log.Printf("Listening on %s", l)
    log.Fatal(http.ListenAndServe(l,  r))
}
