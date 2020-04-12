package pdnsconfig

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Record is an individual DNS resource record
type Record struct {
	Type  string
	Class string
	TTL   int32
	RDATA string
}

// Domain is a struct containing all of the resource records associated with a
// given domain or node.
type Domain struct {
	Records []Record
}

// Configuration is the top-level config struct for an instance of pocketDNS
type Configuration struct {
	Domains map[string]Domain
}

func loadConfig() {
	p, err := os.UserHomeDir()
	if err != nil {
		log.Warn("Unable to get user's home directory, skipping user's config")
	}
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if p != "" {
		//See https://hermanschaaf.com/efficient-string-concatenation-in-go/
		//  for why + is used here instead of strings.Join()
		viper.AddConfigPath(p + ".homeDNS")
	}
	viper.AddConfigPath("/etc/appname/")
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		log.Fatal("Unable to find any config file")
	}
}

// Config takes no arguments and returns a pocketDNS *Configuration if
// possible. On error, will log a Fatal error and quit since running the server
// regardless would result in a DNS server with no records listening only to
// queries from the localhost.
func Config() *Configuration {
	var C Configuration
	loadConfig()
	err := viper.Unmarshal(&C)
	if err != nil {
		log.Fatal("Failed unmarshalling the config: " + err.Error())
	}
	return &C
}
