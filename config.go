package main

type Config struct {
	UseSSL      bool   `json:"use_ssl"`
	SrvPort     string `json:"srv_port"`
	SSLCertFile string `json:"ssl_cert_file"`
	SSLKeyFile  string `json:"ssl_key_file"`
}

func (c *Config) String() string {
	s := "\n------server config------"
	s += "\nssl cert file:" + c.SSLCertFile
	s += "\nssl key file:" + c.SSLKeyFile
	s += "\n-------------------------"
	return s
}

var (
	_sysConfig *Config = nil
)
