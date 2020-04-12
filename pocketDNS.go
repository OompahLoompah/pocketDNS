package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/OompahLoompah/pocketDNS/internal/pdnsconfig"
	dns "github.com/OompahLoompah/pocketDNS/pkg/DNS"
	listener "github.com/OompahLoompah/pocketDNS/pkg/Listener"
)

func parseRecords(domains map[string]pdnsconfig.Domain) map[string]dns.ResourceRecord {
	records := make(map[string]dns.ResourceRecord)
	for n, d := range domains {
		for _, r := range d.Records {
			if r.Type == "A" {
				records[n] = dns.ResourceRecord{
					NAME:     n,
					TYPE:     "A",
					CLASS:    "IN",
					TTL:      uint32(r.TTL),
					RDLENGTH: uint16(len(r.RDATA)),
					RDATA:    r.RDATA,
				}
			}
		}
	}
	return records
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Print("pocketDNS Starting...")

	conf := pdnsconfig.Config()
	for k, v := range conf.Domains {
		log.Debug(k)
		log.Debug(v)
	}
	d := parseRecords(conf.Domains)
	f := &dns.ResponseFactory{
		ARecords: d,
	}
	l := listener.New("127.0.0.1", 53, f)
	err := l.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
