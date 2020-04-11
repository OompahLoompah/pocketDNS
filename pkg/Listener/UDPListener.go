package listener

import (
	"net"

	log "github.com/sirupsen/logrus"

	//TODO: Remove dependency on gopacket
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	dns "github.com/OompahLoompah/pocketDNS/pkg/DNS"
)

// UDPListener defines a UDPListener with IP, Port, and a map of DNS records.
// By default, a UDPListener will use Port 53 but will only listen on the IP
// address 127.0.0.1 to prevent accidentally running on unintended IP
// addresses.
type UDPListener struct {
	IP      string
	Port    int
	Factory *dns.ResponseFactory
}

func (l *UDPListener) respond(b []byte, addr net.Addr, conn *net.UDPConn) {
	packet := gopacket.NewPacket(b, layers.LayerTypeDNS, gopacket.Default)
	dnsPacket := packet.Layer(layers.LayerTypeDNS)
	tcp, _ := dnsPacket.(*layers.DNS)
	answer := l.Factory.BuildResponse(tcp)
	buf := gopacket.NewSerializeBuffer()
	o := gopacket.SerializeOptions{} // See SerializeOptions for more details.
	err := answer.SerializeTo(buf, o)
	if err != nil {
		log.Error("Error writing to buffer") //TODO improve this to handle request tracing
	}
	conn.WriteTo(buf.Bytes(), addr)
}

// Listen takes a list of dns.ResourceRecords to listen for and returns nothing
func (l *UDPListener) Listen() {
	// Get all of our mise en place before we do any network setup
	// For now we only support A records so as a (very) dirty shortcut assume
	//   all records are A records.

	if l.Factory == nil {
		log.Fatal("No response factory provided")
	}
	if l.IP == "" {
		l.IP = "127.0.0.1"
	}
	if l.Port == 0 {
		l.Port = 53
	}

	//Listen on UDP Port
	addr := &net.UDPAddr{
		Port: l.Port,
		IP:   net.ParseIP(l.IP),
	}
	u, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Error(err)
	}

	//TODO: QUESTION: How do golang's contexts work? I know we'll need to support parallel requests and afaict we'll need to use contexts in some way.

	for {
		b := make([]byte, 1024)
		n, cAddr, err := u.ReadFrom(b)
		b = b[:n] // Hack to prevent packet decoding from failing
		if err != nil {
			log.Error(err)
		}
		go l.respond(b, cAddr, u)
	}
}
