package main

import (
	"fmt"

	"github.com/OompahLoompah/pocketDNS/internal/pDNSconfig"
	dns "github.com/OompahLoompah/pocketDNS/pkg/DNSResourceRecord"
	"github.com/OompahLoompah/pocketDNS/pkg/listener"
)

func parseRecords(domains map[string]pDNSconfig.Domain) *map[string]dns.ResourceRecord {
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
	return &records
}

func main() {
	conf, err := pDNSconfig.Config()
	if err != nil {
		fmt.Println(err)
		panic("Unable to parse config")
	}
	d := parseRecords(conf.Domains)
	l := listener.UDPListener{
		IP:      "127.0.0.1",
		Port:    53,
		Records: *d,
	}
	l.Listen()
}
