package dns

//RRType representation of Resource Record types
type RRType string

const (
	//A or IPv4 record
	A RRType = "A"
	//AAAA or IPv6 record
	AAAA RRType = "AAAA"
)

//RRClass is a representation of an RFC1035 DNS resource record class
type RRClass string

const (
	//IN is the prefix code for the INTERNET class
	IN RRClass = "IN"
)

//ResourceRecord is a representation of an RFC1035 DNS Resource Record
type ResourceRecord struct {
	NAME     string
	TYPE     RRType
	CLASS    RRClass
	TTL      uint32
	RDLENGTH uint16
	RDATA    string
}
