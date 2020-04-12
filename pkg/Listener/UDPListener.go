package listener

import (
	"errors"
	"net"

	log "github.com/sirupsen/logrus"

	//TODO: Remove dependency on gopacket
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"

	dns "github.com/OompahLoompah/pocketDNS/pkg/DNS"
)

// UDPListener defines a listener using UDP with IP, Port, and, to handle the
// construction of query responses, a ResponseFactory. By default, a
// UDPListener will use Port 53 but will only listen on the IP address
// 127.0.0.1 to prevent accidentally running on public IP addresses.
type UDPListener struct {
	IP      string
	Port    int
	Factory *dns.ResponseFactory
}

// New builds and returns a new UDPListener. This function enforces the
// existence of required struct members and is considered the safe way
// to create UDPListeners.
func New(ip string, port int, factory *dns.ResponseFactory) *UDPListener {
	return &UDPListener{
		IP:      ip,
		Port:    port,
		Factory: factory,
	}
}

// respond is an internal helper function to handle the construction and
// transmission of answers to individual queries.
func (l *UDPListener) respond(b []byte, addr net.Addr, conn *net.UDPConn) {
	packet := gopacket.NewPacket(b, layers.LayerTypeDNS, gopacket.Default)
	dnsPacket := packet.Layer(layers.LayerTypeDNS)
	tcp, _ := dnsPacket.(*layers.DNS)
	answer := l.Factory.BuildResponse(tcp)
	buf := gopacket.NewSerializeBuffer()
	o := gopacket.SerializeOptions{}
	err := answer.SerializeTo(buf, o)
	if err != nil {
		log.Error("Error writing to buffer") //TODO improve this to handle request tracing
	}
	conn.WriteTo(buf.Bytes(), addr)
}

// Listen handles the initial work of opening the port l.Port on l.IP and, on
// receiving packets, will hand them off to the respond() function as a new
// goroutine. Returns nothing and on error should log it and return the
// appropriate error message to the query sender if possible.
func (l *UDPListener) Listen() error {

	if l.Factory == nil {
		return errors.New("UDPlistener is missing a response factory")
	}
	if l.IP == "" {
		l.IP = "127.0.0.1"
	}
	if l.Port == 0 {
		l.Port = 53
	}

	addr := &net.UDPAddr{
		Port: l.Port,
		IP:   net.ParseIP(l.IP),
	}
	log.Info("Now listening on " + addr.String())
	u, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	for {
		b := make([]byte, 1024)
		n, cAddr, err := u.ReadFrom(b)
		if err != nil {
			log.Error(err)
		}
		b = b[:n] // Trim slice so packet decoding won't fail
		go l.respond(b, cAddr, u)
	}
}
