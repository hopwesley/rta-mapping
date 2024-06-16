package main

import "fmt"

type Config struct {
	UseSSL      bool      `json:"use_ssl"`
	SrvPort     string    `json:"srv_port"`
	SSLCertFile string    `json:"ssl_cert_file"`
	SSLKeyFile  string    `json:"ssl_key_file"`
	RedisCfg    *RedisCfg `json:"redis"`
	MysqlCfg    *MysqlCfg `json:"mysql"`
}

func (c *Config) String() string {
	s := "\n------------server config------------"
	s += fmt.Sprintf("\nif use ssl:%t", c.UseSSL)
	s += "\nserver port:" + c.SrvPort
	s += "\nssl cert file:" + c.SSLCertFile
	s += "\nssl key file:" + c.SSLKeyFile
	s += "\nredis config:" + c.RedisCfg.String()
	s += "\nmysql config:" + c.MysqlCfg.String()
	s += "\n-------------------------------------"
	return s
}

var (
	_sysConfig *Config = nil
)
