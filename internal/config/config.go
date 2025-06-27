package config

import (
	"flag"
	"log"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server Server `mapstructure:"server"`
	Vite   Vite   `mapstructure:"vite"`
}

// Server configuration
type Server struct {
	Addr string `mapstructure:"addr"`
	Dev  bool   `mapstructure:"dev"`
}

// Vite configuration
type Vite struct {
	URL   string `mapstructure:"url"`
	Entry string `mapstructure:"entry"`
}

// Load loads configuration from flags, environment variables, and config files
func Load() *Config {
	// Set up Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set default values
	viper.SetDefault("server.addr", ":8080")
	viper.SetDefault("server.dev", false)
	viper.SetDefault("vite.url", "http://localhost:5173")
	viper.SetDefault("vite.entry", "/src/main.ts")

	// Environment variable bindings
	viper.SetEnvPrefix("MINIO_ADMIN")
	viper.AutomaticEnv()

	// Bind environment variables
	viper.BindEnv("server.addr", "ADDR")
	viper.BindEnv("server.dev", "DEV")
	viper.BindEnv("vite.url", "VITE_URL")
	viper.BindEnv("vite.entry", "VITE_ENTRY")

	// Parse command line flags
	addr := flag.String("addr", viper.GetString("server.addr"), "HTTP server address")
	dev := flag.Bool("dev", viper.GetBool("server.dev"), "Enable development mode")
	flag.Parse()

	// Override with flag values if provided
	if *addr != viper.GetString("server.addr") {
		viper.Set("server.addr", *addr)
	}
	if *dev != viper.GetBool("server.dev") {
		viper.Set("server.dev", *dev)
	}

	// Try to read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Error reading config file: %v", err)
		}
	}

	// Unmarshal configuration
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode configuration: %v", err)
	}

	return &cfg
}
