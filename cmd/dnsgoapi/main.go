package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"dnsgoapi/dnsgoapi"
)

func main() {
	r := mux.NewRouter()	

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to dnsgoapi! Look at GitHub for how to interact with this API!")
	}).Methods("GET")

	r.HandleFunc("/a/{PublicDNS}/{q}", dnsgoapi.QueryA).Methods("GET")
	r.HandleFunc("/aaaa/{PublicDNS}/{q}", dnsgoapi.QueryQuadA).Methods("GET")
	r.HandleFunc("/cname/{PublicDNS}/{q}", dnsgoapi.QueryCNAME).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}