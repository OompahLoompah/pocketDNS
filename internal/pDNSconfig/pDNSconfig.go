package pDNSconfig

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Record struct {
	Type  string
	Class string
	TTL   int32
	RDATA string
}

type Domain struct {
	Records []Record
}

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

var C Configuration

func Config() *Configuration {
	loadConfig()
	err := viper.Unmarshal(&C)
	if err != nil {
		log.Fatal("Failed unmarshalling the config: " + err.Error())
	}
	return &C
}
