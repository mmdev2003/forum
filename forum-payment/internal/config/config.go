package config

func New() *Config {
	return &Config{
		DB: DB{
			Username: "postgres",
			Dbname:   "postgres",
			Password: "postgres",
			Host:     "forum-payment-db",
			Port:     "5432",
		},
		Forum: Forum{
			Authorization: Server{
				Host: "forum-authorization",
				Port: "8001",
			},
			Frame: Server{
				Host: "forum-frame",
				Port: "8004",
			},
			Status: Server{
				Host: "forum-status",
				Port: "8005",
			},
			Payment: Server{
				Host: "forum-payment",
				Port: "8009",
			},
			InterServerSecretKey: "secret-key",
			BtcAddress:           "1DdiKnSd4pZ69S4oRWod66e6G3hWaSZXAV",
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
	Frame                Server
	InterServerSecretKey string
	BtcAddress           string
}

type Server struct {
	Host string
	Port string
}
