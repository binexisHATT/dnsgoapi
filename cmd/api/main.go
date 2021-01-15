package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/dnsgoapi/pkg/dnsapi"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
)

func index(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("./README.md")
	if err != nil {
		log.Fatal("Unable to open README.md file")
	}
	defer f.Close()

	markdown, err := ioutil.ReadAll(f)

    fmt.Fprintf(w, string(blackfriday.MarkdownCommon(markdown)))
}

func main() {
	port := flag.Int("port", 8080, "The port for the server to listen on")
	flag.Parse()

	r := mux.NewRouter()

    r.HandleFunc("/", index)

    r.HandleFunc(
        "/{recordType:[a-zA-Z]+}/{publicDNS:[a-zA-Z9]+}/{fqdn}",
        dnsgoapi.DNSQuery,
    ).Methods("GET")

    http.Handle("/", r)

	l := fmt.Sprintf(":%d", *port)
	log.Printf("Listening on %s", l)

	log.Fatal(http.ListenAndServe(l, r))
}
