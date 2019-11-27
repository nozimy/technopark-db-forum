package apiserver

type Config struct {
	BindAddr    string
	SessionKey  string
	DatabaseURL string
}

func NewConfig() *Config {
	return &Config{
		BindAddr:    ":5000",
		SessionKey:  "jdfhdfdj",
		DatabaseURL: "host=localhost dbname=docker sslmode=disable port=5432 password=docker user=docker",
	}
}
