package config

func New() *Config {
	return &Config{
		DB: DB{
			Username: "postgres",
			Dbname:   "postgres",
			Password: "postgres",
			Host:     "forum-dialog-db",
			Port:     "5432",
		},
		Forum: Forum{
			Authorization: Server{
				Host: "forum-authorization",
				Port: "8001",
			},
			Notification: Server{
				Host: "forum-notification",
				Port: "8002",
			},
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
	Notification  Server
}

type Server struct {
	Host string
	Port string
}
