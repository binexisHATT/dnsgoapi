package dnsgoapi

import (
	"encoding/json"
	"log"
	"net/http"
	re "regexp"

	"github.com/gorilla/mux"
	"github.com/miekg/dns"
)

var (
	msg dns.Msg
)

func matchString(expression, serverName string) bool {
	matched, _ := re.MatchString(expression, serverName)
	return matched
}

func getDNSIP(s string) string {
	switch {
	case matchString("(?i)cloudflare", s):
		return "1.1.1.1:53"
	case matchString("(?i)google", s):
		return "8.8.8.8:53"
	case matchString("(?i)opendns", s):
		return "208.67.222.222:53"
	case matchString("(?i)comodo", s):
		return "8.26.56.26:53"
	case matchString("(?i)quad9", s):
		return "9.9.9.9:53"
	case matchString("(?i)verisign", s):
		return "64.6.64.6:53"
	default:
		return "1.1.1.1:53"
	}
}

func setRecordType(fqdn, requestedRecord string, msg *dns.Msg) {
	fqdn = dns.Fqdn(fqdn)
	switch {
	case matchString("(?i)^a$", requestedRecord):
		log.Println("(match) A")
		msg.SetQuestion(fqdn, dns.TypeA)
	case matchString("(?i)aaaa", requestedRecord):
		log.Println("(match) AAAA")
		msg.SetQuestion(fqdn, dns.TypeAAAA)
	case matchString("(?i)cname", requestedRecord):
		log.Println("(match) CNAME")
		msg.SetQuestion(fqdn, dns.TypeCNAME)
	case matchString("(?i)mx", requestedRecord):
		log.Println("(match) MX")
		msg.SetQuestion(fqdn, dns.TypeMX)
	case matchString("(?i)ns", requestedRecord):
		log.Println("(match) NS")
		msg.SetQuestion(fqdn, dns.TypeNS)
	case matchString("(?i)txt", requestedRecord):
		log.Println("(match) TXT")
		msg.SetQuestion(fqdn, dns.TypeTXT)
	default:
		msg.SetQuestion(fqdn, dns.TypeA)
	}
}

func DNSQuery(w http.ResponseWriter, r *http.Request) {
	result := make(map[string][]string)

	requestedRecord := mux.Vars(r)["recordType"]
	fqdn := mux.Vars(r)["fqdn"]
	publicDNS := mux.Vars(r)["publicDNS"]

	DNSServerIP := getDNSIP(publicDNS)
	setRecordType(fqdn, requestedRecord, &msg)

	var (
		resp *dns.Msg
		err  error
	)

	resp, err = dns.Exchange(&msg, DNSServerIP)
	if err != nil {
		http.Redirect(w, r, "/", 302)
	}

	if len(resp.Answer) < 1 {
		log.Println("No answers")
		result[fqdn] = append(result[fqdn], "No answers")
	} else {
		for _, ans := range resp.Answer {
			switch ans.(type) {
			case *dns.A:
				result[fqdn] = append(result[fqdn], ans.(*dns.A).A.String())
			case *dns.AAAA:
				result[fqdn] = append(result[fqdn], ans.(*dns.AAAA).AAAA.String())
			case *dns.CNAME:
				result[fqdn] = append(result[fqdn], ans.(*dns.CNAME).Target)
			case *dns.MX:
				result[fqdn] = append(result[fqdn], ans.(*dns.MX).Mx)
			case *dns.NS:
				result[fqdn] = append(result[fqdn], ans.(*dns.NS).Ns)
			case *dns.TXT:
				result[fqdn] = append(result[fqdn], ans.(*dns.TXT).String())
			default:
				result[fqdn] = append(result[fqdn], "No answers")
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
}
