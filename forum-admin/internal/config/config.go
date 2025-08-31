package config

func New() *Config {
	return &Config{
		DB: DB{
			Username: "postgres",
			Dbname:   "postgres",
			Password: "postgres",
			Host:     "forum-admin-db",
			Port:     "5432",
		},
		Forum: Forum{
			Admin: Server{
				Host: "forum-admin",
				Port: "8010",
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
	Admin                Server
	InterServerSecretKey string
}

type Server struct {
	Host string
	Port string
}
