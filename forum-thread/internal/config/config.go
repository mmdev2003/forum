package config

func New() *Config {
	return &Config{
		DB: DB{
			Username: "postgres",
			Dbname:   "postgres",
			Password: "postgres",
			Host:     "forum-thread-db",
			Port:     "5432",
		},
		Meilisearch: Meilisearch{
			"forum-meilisearch",
			"7700",
			"qwertyuiop",
		},
		Rabbitmq: RabbitMQ{
			"forum-rabbitmq",
			"5672",
			"aloha",
			"tool_pool",
		},
		WeedFS: WeedFS{
			FilerUrl:  "http://weed-filer:8888",
			MasterUrl: "http://weed-master:9333",
		},
		Forum: Forum{
			Authorization: Server{
				Host: "forum-authorization",
				Port: "8001",
			},
			Thread: Server{
				Host: "forum-thread",
				Port: "8003",
			},
			Status: Server{
				Host: "forum-status",
				Port: "8005",
			},
			User: Server{
				Host: "forum-user",
				Port: "8006",
			},
			Notification: Server{
				Host: "forum-notification",
				Port: "8007",
			},
			Support: Server{
				Host: "forum-support",
				Port: "8008",
			},
			InterServerSecretKey: "inter_server_secret_key",
		},
	}
}

type Config struct {
	DB          DB
	Meilisearch Meilisearch
	Rabbitmq    RabbitMQ
	WeedFS      WeedFS
	Forum       Forum
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
	Support              Server
	Status               Server
	Thread               Server
	Notification         Server
	User                 Server
	InterServerSecretKey string
}

type Meilisearch struct {
	Host   string
	Port   string
	ApiKey string
}

type WeedFS struct {
	FilerUrl  string
	MasterUrl string
}

type RabbitMQ struct {
	Host     string
	Port     string
	Username string
	Password string
}

type Server struct {
	Host string
	Port string
}
