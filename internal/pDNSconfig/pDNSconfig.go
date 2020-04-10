package pDNSconfig

import (
	"fmt"
	"os"

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
		fmt.Println("Unable to get user's home directory, unable to parse any ",
			"config file in $HOME/.homeDNS") //TODO: send this to logger instead of Println
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
		fmt.Println("Unable to find config file")
	}
}

var C Configuration

func Config() (*Configuration, error) {
	loadConfig()
	err := viper.Unmarshal(&C)
	if err != nil {
		return nil, err
	}
	return &C, nil
}
