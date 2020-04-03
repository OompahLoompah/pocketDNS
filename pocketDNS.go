package main

import (
	//    "net"

	//    "github.com/google/gopacket"
	//    layers "github.com/google/gopacket/layers"

	"fmt"

	c "github.com/OompahLoompah/pocketDNS/internal/config"
	dns "github.com/OompahLoompah/pocketDNS/pkg/DNSResourceRecord"
)

type config struct {
	Domains []dns.ResourceRecord
}

func main() {
	conf, err := c.Config()
	if err != nil {
		fmt.Println(err)
		panic("Unable to parse config")
	}
	fmt.Println((*conf).Domains["home.lab"].Records[0])
}
