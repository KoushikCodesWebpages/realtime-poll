package config

type Config struct {
	AppEnv   string
	Debug    bool
	Port     string
	MongoURI string
	DBName   string
}

var App Config

func Defaults() {
	App = Config{
		AppEnv:   "dev",
		Debug:    true,
		Port:     "8080",
		MongoURI: "mongodb://localhost:27017",
		DBName:   "realtime_poll_dev",
	}
}
