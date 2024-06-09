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
	}
	bts, _ := json.MarshalIndent(cfg, "", "\t")
	_ = os.WriteFile("../config.sample.json", bts, 0644)
}
