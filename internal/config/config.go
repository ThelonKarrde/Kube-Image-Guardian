package config

import (
	"os"
)

type StartUpConfig struct {
	TlsCertPath string
	TlsKeyPath  string
	Port        string
	ConfigPath  string
}

func (c *StartUpConfig) InitConfig() {
	var ok bool
	c.ConfigPath, ok = os.LookupEnv("CONFIG_PATH")
	if !ok {
		c.ConfigPath = "/etc/app-config/config/config.yaml"
	}
	c.TlsCertPath, ok = os.LookupEnv("TLS_CERT_PATH")
	if !ok {
		c.TlsCertPath = "/etc/app-config/tls/tls.crt"
	}
	c.TlsKeyPath, ok = os.LookupEnv("TLS_KEY_PATH")
	if !ok {
		c.TlsKeyPath = "/etc/app-config/tls/tls.key"
	}
	c.Port, ok = os.LookupEnv("PORT")
	if !ok {
		c.Port = "1224"
	}
}
