package config

func New() *Config {
	return &Config{
		DB: DB{
			Username: "postgres",
			Dbname:   "postgres",
			Password: "postgres",
			Host:     "forum-authentication-db",
			Port:     "5432",
		},
		Forum: Forum{
			Authentication: Server{
				Host: "forum-authentication",
				Port: "8000",
			},
			Thread: Server{
				Host: "forum-thread",
				Port: "8002",
			},
			Authorization: Server{
				Host: "forum-authorization",
				Port: "8001",
			},
			User: Server{
				Host: "forum-user",
				Port: "8009",
			},
			Admin: Server{
				Host: "forum-admin",
				Port: "8010",
			},
			InterServerSecretKey: "inter_server_secret_key",
			PasswordSecretKey:    "password_secret_key",
			Domain:               "212.113.122.80",
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
	Authentication       Server
	Thread               Server
	User                 Server
	Admin                Server
	InterServerSecretKey string
	PasswordSecretKey    string
	Domain               string
}

type Server struct {
	Host string
	Port string
}
