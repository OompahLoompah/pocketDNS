package dns

import (
	"net"

	log "github.com/sirupsen/logrus"

	//TODO: Remove dependency on gopacket
	"github.com/google/gopacket/layers"
)

// ResponseFactory Config
type ResponseFactory struct {
	ARecords map[string]ResourceRecord
}

//BuildResponse builds and returns a DNS answer.
func (d *ResponseFactory) BuildResponse(request *layers.DNS) *layers.DNS {
	replyMess := request // Using the request as a starter for the response

	//Additionals seems to sometimes have an emmpty element. We need to reset
	//the additionals section as a result
	replyMess.Additionals = nil
	replyMess.ARCount = 0

	for _, q := range request.Questions {
		switch q.Type {
		case layers.DNSTypeA:
			record, ok := d.ARecords[string(request.Questions[0].Name)]
			if !ok {
				log.Debug("Got request for A record of unknown domain: " + string(request.Questions[0].Name))
			} else {
				addr := net.ParseIP(record.RDATA)
				dnsAnswer := layers.DNSResourceRecord{
					Name:  request.Questions[0].Name,
					Type:  layers.DNSTypeA,
					Class: layers.DNSClassIN,
					TTL:   record.TTL,
					IP:    addr,
				}
				replyMess.Answers = append(replyMess.Answers, dnsAnswer)
				replyMess.ANCount++
			}
		}
	}
	replyMess.QR = true // message is a response, not a query
	replyMess.ResponseCode = layers.DNSResponseCodeNoErr
	return replyMess
}
