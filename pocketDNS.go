package main

import (
	"fmt"

	c "github.com/OompahLoompah/pocketDNS/internal/config"
	dns "github.com/OompahLoompah/pocketDNS/pkg/DNSResourceRecord"
	"github.com/OompahLoompah/pocketDNS/pkg/listener"
)

type config struct {
	Domains []dns.ResourceRecord
}

func parseRecords(domains map[string]c.Domain) *map[string]dns.ResourceRecord {
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
	conf, err := c.Config()
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
