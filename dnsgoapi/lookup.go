package dnsgoapi

import (
	"encoding/json"
	"net/http"
	re "regexp"

	"github.com/miekg/dns"
	"github.com/gorilla/mux"
)

var (
	msg dns.Msg
)

func getPublicDNSServer(s string) string {
	switch s {
	case re.MatchString("(?i)cloudflare"):
		return "1.1.1.1:53"
	case re.MatchString("(?i)google"):
		return "8.8.8.8:53"
	case re.MatchString("(?i)opendns"):
		return "208.67.222.222:53"
	case re.MatchString("(?i)comodo"):
		return "8.26.56.26:53"	
	case re.MatchString("(?i)quad9"):
		return "9.9.9.9:53"
	case re.MatchString("(?i)verisign"):
		return "64.6.64.6:53"
	default:
		return "1.1.1.1:53"
	}
}

func QueryA(w http.ResponseWriter, r *http.Request) {
	result := make(map[string][]string)

	q := mux.Vars(r)["q"]
	publicDNS := mux.Vars(r)["DNSServer"]

	publicDNS = getPublicDNSServer(publicDNS)

	fqdn := dns.Fqdn(q)
	msg.SetQuestion(fqdn, dns.TypeA)
	resp, err := dns.Exchange(&msg, publicDNS)
	if err != nil {
		panic(err)
	}

	if len(resp.Answer) < 1 {
		result[fqdn] = append(result[fqdn], "No answers")
	}

	for _, answer := range resp.Answer {
		if a, ok := answer.(*dns.A); ok {
			result[fqdn] = append(result[fqdn], a.A.String())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
}

func QueryQuadA(w http.ResponseWriter, r *http.Request) {
	//result := make(map[string][]string)

	q := mux.Vars(r)["q"]
	//publicDNS := mux.Vars(r)["DNSServer"]
	fqdn := dns.Fqdn(q)
	msg.SetQuestion(fqdn, dns.TypeAAAA)

	w.Header().Set("Content-Type", "application/json")
}

func QueryCNAME(w http.ResponseWriter, r *http.Request) {
	//result := make(map[string][]string)

	q := mux.Vars(r)["q"]
	//publicDNS := mux.Vars(r)["DNSServer"]
	fqdn := dns.Fqdn(q)
	msg.SetQuestion(fqdn, dns.TypeCNAME)


	w.Header().Set("Content-Type", "application/json")
}

func QueryTXT(w http.ResponseWriter, r *http.Request) {
	//result := make(map[string][]string)

	q := mux.Vars(r)["q"]
	//publicDNS := mux.Vars(r)["DNSServer"]
	fqdn := dns.Fqdn(q)
	msg.SetQuestion(fqdn, dns.TypeTXT)

	w.Header().Set("Content-Type", "application/json")
}

func QueryMX(w http.ResponseWriter, r *http.Request) {
	//result := make(map[string][]string)

	q := mux.Vars(r)["q"]
	//publicDNS := mux.Vars(r)["DNSServer"]
	fqdn := dns.Fqdn(q)
	msg.SetQuestion(fqdn, dns.TypeMX)

	w.Header().Set("Content-Type", "application/json")
}
