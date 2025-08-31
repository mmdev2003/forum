package config

func New() *Config {
	return &Config{
		DB: DB{
			Username: "postgres",
			Dbname:   "postgres",
			Password: "postgres",
			Host:     "forum-authorization-db",
			Port:     "5432",
		},
		Forum: Forum{
			Authorization: Server{
				Host: "forum-authorization",
				Port: "8001",
			},
			JwtSecretKey: "jwt_secret_key",
			Domain:       "212.113.122.80",
		},
	}
}

type Config struct {
	DB    DB
	Forum Forum
}

type DB struct {
	Username string
	Password string
	Host     string
	Port     string
	Dbname   string
}

type Forum struct {
	Authorization Server
	JwtSecretKey  string
	Domain        string
}

type Server struct {
	Host string
	Port string
}
