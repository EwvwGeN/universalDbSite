package server

import "os"

type Config struct {
	Api_host string `yaml:"api_host"`
	Api_port string `yaml:"api_port"`
	Db_host  string `yaml:"db_host"`
	Db_port  string `yaml:"db_port"`
	Db_user  string `yaml:"db_user"`
	Db_pass  string `yaml:"db_pass"`
	Db_name  string `yaml:"db_name"`
}

func NewConfig() *Config {
	return &Config{
		Api_host: getEnv("API_HOST", "0.0.0.0"),
		Api_port: getEnv("API_PORT", "9090"),
		Db_host:  getEnv("DB_HOST", "0.0.0.0"),
		Db_port:  getEnv("DB_PORT", "3306"),
		Db_user:  getEnv("DB_USER", ""),
		Db_pass:  getEnv("DB_PASS", ""),
		Db_name:  getEnv("DB_NAME", ""),
	}

}

func getEnv(envKey string, defaultVal string) string {
	if value, exists := os.LookupEnv(envKey); exists {
		return value
	}

	return defaultVal
}
