package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"

	"dnsgoapi/dnsgoapi"
)

func main() {
	port := flag.Int("port", 8080, "The port for the server to listen on")
	flag.Parse()

	r := mux.NewRouter()
	f, err := os.Open("./README.md")
	if err != nil {
		log.Fatal("Unable to open README.md file")
	}
	defer f.Close()
	markdown, err := ioutil.ReadAll(f)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(blackfriday.MarkdownCommon(markdown)))
	}).Methods("GET")

	r.HandleFunc("/a/{PublicDNS}/{q}", dnsgoapi.QueryA).Methods("GET")
	r.HandleFunc("/aaaa/{PublicDNS}/{q}", dnsgoapi.QueryQuadA).Methods("GET")
	r.HandleFunc("/cname/{PublicDNS}/{q}", dnsgoapi.QueryCNAME).Methods("GET")

	l := fmt.Sprintf(":%d", *port)
	log.Printf("Listening on %s", l)
	log.Fatal(http.ListenAndServe(l, r))
}
