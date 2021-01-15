package dnsgoapi

import (
	"log"
	"encoding/json"
	"net/http"
	re "regexp"
	"strings"

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

func replaceTabs(str, replace_with string) string {

	return strings.Replace(str, "\t", replace_with, -1)
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
		log.Println("(Match) A Record")
		msg.SetQuestion(fqdn, dns.TypeA)
	case matchString("(?i)aaaa", requestedRecord):
		log.Println("(Match) AAAA Record")
		msg.SetQuestion(fqdn, dns.TypeAAAA)
	case matchString("(?i)cname", requestedRecord):
		log.Println("(Match) CNAME Record")
		msg.SetQuestion(fqdn, dns.TypeCNAME)
	case matchString("(?i)mx", requestedRecord):
		log.Println("(Match) MX Record")
		msg.SetQuestion(fqdn, dns.TypeMX)
	case matchString("(?i)ns", requestedRecord):
		log.Println("(Match) NS Record")
		msg.SetQuestion(fqdn, dns.TypeNS)
	case matchString("(?i)txt", requestedRecord):
		log.Println("(Match) TXT Record")
		msg.SetQuestion(fqdn, dns.TypeTXT)
	case matchString("(?i)ptr", requestedRecord):
		log.Println("(Match) PTR Record")
		msg.SetQuestion(fqdn, dns.TypePTR)
	case matchString("(?i)caa", requestedRecord):
		log.Println("(Match) CAA Record")
		msg.SetQuestion(fqdn, dns.TypeCAA)
	case matchString("(?i)soa", requestedRecord):
		log.Println("(Match) SOA Record")
		msg.SetQuestion(fqdn, dns.TypeSOA)
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
				log.Printf("(Response) A Record Answer for %s", fqdn)
				result[fqdn] = append(result[fqdn], ans.(*dns.A).A.String())
			case *dns.AAAA:
				log.Printf("(Response) AAAA Record Answer for %s", fqdn)
				result[fqdn] = append(result[fqdn], ans.(*dns.AAAA).AAAA.String())
			case *dns.CNAME:
				log.Printf("(Response) CNAME Record Answer for %s", fqdn)
				result[fqdn] = append(result[fqdn], ans.(*dns.CNAME).Target)
			case *dns.MX:
				log.Printf("(Response) MX Record Answer for %s", fqdn)
				result[fqdn] = append(result[fqdn], ans.(*dns.MX).Mx)
			case *dns.NS:
				log.Printf("(Response) NS Record Answer for %s", fqdn)
				result[fqdn] = append(result[fqdn], ans.(*dns.NS).Ns)
			case *dns.TXT:
				log.Printf("(Response) TXT Record Answer for %s", fqdn)
                txt := ans.(*dns.TXT).Txt[0]
				result[fqdn] = append(result[fqdn], txt)
			case *dns.PTR:
				log.Printf("(Response) PTR Record Answer for %s", fqdn)
				result[fqdn] = append(result[fqdn], ans.(*dns.PTR).String())
			case *dns.CAA:
				log.Printf("(Response) CAA Record Answer for %s", fqdn)
				result[fqdn] = append(result[fqdn], ans.(*dns.CAA).Value)
			case *dns.SOA:
				log.Printf("(Response) SOA Record Answer for %s", fqdn)
				// Extracting only the primary name server and primary email
				var temp []string
				a := ans.(*dns.SOA)
				temp = append(temp, a.Ns)
				temp = append(temp, a.Mbox)
				soa := strings.Join(temp, ",")
				result[fqdn] = append(result[fqdn], soa)
			default:
				result[fqdn] = append(result[fqdn], "No answers")
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
}
