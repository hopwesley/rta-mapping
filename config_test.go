package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	cfg := &Config{
		UseSSL:      false,
		SrvPort:     "8801",
		SSLCertFile: "",
		SSLKeyFile:  "",
		Redis: &RedisCfg{
			Addr:     "localhost:6379",
			Password: "123",
		},
	}
	bts, _ := json.MarshalIndent(cfg, "", "\t")
	_ = os.WriteFile("../config.sample.json", bts, 0644)
}
