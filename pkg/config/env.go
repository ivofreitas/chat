package config

import (
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"sync"
)

// Env values
type Env struct {
	Server   Server
	Log      Log
	Doc      Doc
	Database Database
	Broker   Broker
	External External
	Security Security
}

// Server config
type Server struct {
	AuthPort string
	ChatPort string
}

// Log config
type Log struct {
	Enabled bool
	Level   string
}

// Doc - swagger information
type Doc struct {
	Title       string
	Description string
	Enabled     bool
	Version     string
}

// Database - Postgres configuration
type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Broker struct {
	Url             string
	User            string
	Password        string
	StockCodeQueue  string
	StockQuoteQueue string
}

type External struct {
	StockBaseUrl string
}

type Security struct {
	JWTSecretKey string
}

var (
	env  *Env
	once sync.Once
)

// GetEnv returns env values
func GetEnv() *Env {

	once.Do(func() {

		viper.AutomaticEnv()
		err := godotenv.Load("./pkg/config/.env")
		if err != nil {
			log.Warn(err)
		}

		env = new(Env)
		env.Server.AuthPort = viper.GetString("SERVER_AUTH_PORT")
		env.Server.ChatPort = viper.GetString("SERVER_CHAT_PORT")

		env.Log.Enabled = viper.GetBool("LOG_ENABLED")
		env.Log.Level = viper.GetString("LOG_LEVEL")

		env.Database.Host = viper.GetString("DB_HOST")
		env.Database.Port = viper.GetString("DB_PORT")
		env.Database.User = viper.GetString("DB_USER")
		env.Database.Password = viper.GetString("DB_PASSWORD")
		env.Database.DBName = viper.GetString("DB_NAME")
		env.Database.SSLMode = viper.GetString("DB_SSLMODE")

		env.Broker.Url = viper.GetString("BROKER_URL")
		env.Broker.User = viper.GetString("BROKER_USER")
		env.Broker.Password = viper.GetString("BROKER_PASSWORD")
		env.Broker.StockCodeQueue = viper.GetString("BROKER_STOCK_CODE_QUEUE")
		env.Broker.StockQuoteQueue = viper.GetString("BROKER_STOCK_QUOTE_QUEUE")

		env.External.StockBaseUrl = viper.GetString("EXTERNAL_STOCK_BASE_URL")

		env.Security.JWTSecretKey = viper.GetString("SECURITY_JWT_SECRET_KEY")
	})

	return env
}
