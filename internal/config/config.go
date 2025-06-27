package config

import (
	"flag"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server Server `mapstructure:"server"`
	Vite   Vite   `mapstructure:"vite"`
	Logger Logger `mapstructure:"logger"`
	MinIO  MinIO  `mapstructure:"minio"`
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

// Logger configuration
type Logger struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
}

// MinIO configuration
type MinIO struct {
	URL      string `mapstructure:"url"`
	RootUser string `mapstructure:"root_user"`
	Password string `mapstructure:"password"`
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
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.pretty", true)
	viper.SetDefault("minio.url", "http://localhost:9000")
	viper.SetDefault("minio.root_user", "")
	viper.SetDefault("minio.password", "")

	// Environment variable bindings
	viper.SetEnvPrefix("MINIO_ADMIN")
	viper.AutomaticEnv()

	// Bind environment variables
	if err := viper.BindEnv("server.addr", "ADDR"); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind server.addr environment variable")
	}
	if err := viper.BindEnv("server.dev", "DEV"); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind server.dev environment variable")
	}
	if err := viper.BindEnv("vite.url", "VITE_URL"); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind vite.url environment variable")
	}
	if err := viper.BindEnv("vite.entry", "VITE_ENTRY"); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind vite.entry environment variable")
	}
	if err := viper.BindEnv("logger.level", "LOG_LEVEL"); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind logger.level environment variable")
	}
	if err := viper.BindEnv("logger.pretty", "LOG_PRETTY"); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind logger.pretty environment variable")
	}
	if err := viper.BindEnv("minio.url", "MINIO_URL"); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind minio.url environment variable")
	}
	if err := viper.BindEnv("minio.root_user", "MINIO_ROOT_USER"); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind minio.root_user environment variable")
	}
	if err := viper.BindEnv("minio.password", "MINIO_ROOT_PASSWORD"); err != nil {
		log.Fatal().Err(err).Msg("Failed to bind minio.password environment variable")
	}

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
			log.Warn().Err(err).Msg("Error reading config file")
		}
	}

	// Unmarshal configuration
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Unable to decode configuration")
	}

	return &cfg
}
