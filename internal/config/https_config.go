package config

type HTTPSConfig struct {
	Enable   bool   `mapstructure:"ENABLE_HTTPS"`
	CertFile string `mapstructure:"HTTPS_CERT_FILE"`
	KeyFile  string `mapstructure:"HTTPS_KEY_FILE"`
	Port     string `mapstructure:"HTTPS_PORT"`
}
