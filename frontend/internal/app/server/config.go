package server

import "os"

type Config struct {
	Site_host string `yaml:"site_host"`
	Site_port string `yaml:"site_port"`
	Api_host  string `yaml:"api_host"`
	Api_port  string `yaml:"api_port"`
}

func NewConfig() *Config {
	return &Config{
		Site_host: getEnv("SITE_HOST", "0.0.0.0"),
		Site_port: getEnv("SITE_PORT", "9091"),
		Api_host:  getEnv("API_HOST", "localhost"),
		Api_port:  getEnv("API_PORT", "9090"),
	}

}

func getEnv(envKey string, defaultVal string) string {
	if value, exists := os.LookupEnv(envKey); exists {
		return value
	}

	return defaultVal
}
