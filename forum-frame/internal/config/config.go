package config

func New() *Config {
	return &Config{
		DB: DB{
			Username: "postgres",
			Dbname:   "postgres",
			Password: "postgres",
			Host:     "forum-frame-db",
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
			Payment: Server{
				Host: "forum-payment",
				Port: "8009",
			},
			InterServerSecretKey: "inter_server_secret_key",
		},
		WeedFS: WeedFS{
			FilerUrl:  "http://weed-filer:8888",
			MasterUrl: "http://weed-master:9333",
		},
	}
}

type Config struct {
	DB     DB
	WeedFS WeedFS
	Forum  Forum
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
	Frame                Server
	InterServerSecretKey string
}

type Server struct {
	Host string
	Port string
}

type WeedFS struct {
	FilerUrl  string
	MasterUrl string
}
