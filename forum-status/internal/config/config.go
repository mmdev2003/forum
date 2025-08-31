package config

func New() *Config {
	return &Config{
		DB: DB{
			Username: "postgres",
			Dbname:   "postgres",
			Password: "postgres",
			Host:     "forum-status-db",
			Port:     "5432",
		},
		Forum: Forum{
			Authorization: Server{
				Host: "forum-authorization",
				Port: "8001",
			},
			Payment: Server{
				Host: "forum-payment",
				Port: "8009",
			},
			Status: Server{
				Host: "forum-status",
				Port: "8005",
			},
			InterServerSecretKey: "inter_server_secret_key",
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
	Authorization        Server
	Payment              Server
	Status               Server
	InterServerSecretKey string
}

type Server struct {
	Host string
	Port string
}
