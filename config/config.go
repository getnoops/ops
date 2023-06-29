package config

type Auth struct {
	Host     string
	ClientId string
	Scopes   []string
}

type Config struct {
	Auth Auth
}

func New() *Config {
	return &Config{
		Auth: Auth{
			Host:     "github.com",
			ClientId: "05a4d23d91dced4130ce",
		},
	}
}
