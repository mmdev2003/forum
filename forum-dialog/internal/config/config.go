package config

import "github.com/google/uuid"

func New() *Config {
	return &Config{
		Redis: Redis{
			Host:     "forum-redis",
			Port:     "6379",
			Password: "password111",
		},
		RabbitMQ: RabbitMQ{
			Host:     "forum-rabbitmq",
			Port:     "5672",
			Username: "aloha",
			Password: "tool_pool",
		},
		DB: DB{
			Username: "postgres",
			Dbname:   "postgres",
			Password: "postgres",
			Host:     "forum-dialog-db",
			Port:     "5432",
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
			Dialog: Server{
				Host: "forum-dialog",
				Port: "8002",
			},
		},
		PodID: uuid.NewString(),
	}
}

type Config struct {
	Redis    Redis
	RabbitMQ RabbitMQ
	DB       DB
	Forum    Forum
	WeedFS   WeedFS
	PodID    string
}

type RabbitMQ struct {
	Host     string
	Port     string
	Username string
	Password string
}

type DB struct {
	Username string
	Password string
	Host     string
	Port     string
	Dbname   string
}

type Redis struct {
	Host     string
	Port     string
	Password string
}

type WeedFS struct {
	FilerUrl  string
	MasterUrl string
}

type Forum struct {
	Authorization Server
	Dialog        Server
}

type Server struct {
	Host string
	Port string
}
