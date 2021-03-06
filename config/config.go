package config

import (
	"fmt"

	"github.com/exlinc/golang-utils/envconfig"
	"github.com/sirupsen/logrus"
)

// Config - The envconfig struct tag is used to explicitly name the var, set defaults, and flag required values
type Config struct {
	DBPath             string   `envconfig:"DB_PATH" required:"true"`
	Mode               string   `envconfig:"MODE" default:"production"`
	ListenAddress      string   `envconfig:"LISTEN_ADDRESS" default:"0.0.0.0"`
	ListenPort         string   `envconfig:"LISTEN_PORT" default:"3333"`
	AllowedOrigins     []string `envconfig:"ALLOWED_ORIGINS" default:"*"`
	ServiceAPIKey      string   `envconfig:"SERVICE_API_KEY" default:"insecure"`
	TokenCookieName    string   `envconfig:"TOKEN_COOKIE_NAME" default:"auth_tkn"`
	AccessSecret       string   `envconfig:"ACCESS_SECRET" default:"ACCESS_SECRET_KEY"`
	SMTPEmail          string   `envconfig:"SMTP_EMAIL" default:"admin@hiottech.com"`
	SMTPPassword       string   `envconfig:"SMTP_PASSWORD" default:"sx7zhqCSMvfy"`
	SMTPHost           string   `envconfig:"SMTP_HOST" default:"smtp.zoho.com"`
	SMTPPort           string   `envconfig:"SMTP_PORT" default:"587"`
	StorageBucket      string   `envconfig:"STORAGE_BUCKET" default:"file-bucket"`
	StorageKeyLocation string   `envconfig:"STORAGE_KEY_Location" default:"./file-key.json"`
	APIURL             string   `envconfig:"API_URL" default:"http://api.example.com/"`
}

var conf *Config

const (
	// DebugMode - Debug Mode Definition
	DebugMode = "debug"
	// ProductionMode - Production Mode definition
	ProductionMode = "production"
)

// This function gets called automatically when the package is loaded
func init() {
	conf = &Config{}
	// This prefix means our variables will be in the form of GBASE_MODE (for example)
	err := envconfig.Process("GBASE", conf)
	if err != nil {
		fmt.Println("Fatal error processing configuration")
		panic(err)
	}
	l := conf.GetLogger()

	// Sanity check
	if !conf.IsDebugMode() && !conf.IsProductionMode() {
		l.Fatal("Invalid GBASE_MODE variable, it must be either `debug` or `production`")
	}
}

// Cfg returns the configuration - will panic if the config has not been loaded or is nil (which shouldn't happen as that's implicit in the package init)
func Cfg() *Config {
	if conf == nil {
		panic("Config is nil")
	}
	return conf
}

// GetLogger Get Logger
func (cfg *Config) GetLogger() *logrus.Logger {
	var l = logrus.New()
	l.Formatter = &logrus.JSONFormatter{}
	return l
}

// IsDebugMode - returns if Debug Mode
func (cfg *Config) IsDebugMode() bool {
	return cfg.Mode == DebugMode
}

// IsProductionMode - returns if Production Mode
func (cfg *Config) IsProductionMode() bool {
	return cfg.Mode == ProductionMode
}
