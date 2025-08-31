package config

import "github.com/google/uuid"

func New() *Config {
	return &Config{
		Redis: Redis{
			Host:     "forum-redis",
			Port:     "6379",
			Password: "password111",
		},
		DB: DB{
			Username: "postgres",
			Dbname:   "postgres",
			Password: "postgres",
			Host:     "forum-support-db",
			Port:     "5432",
		},
		Rabbitmq: RabbitMQ{
			"forum-rabbitmq",
			"5672",
			"aloha",
			"tool_pool",
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
			Support: Server{
				Host: "forum-support",
				Port: "8002",
			},
		},
		PodID: uuid.NewString(),
	}
}

type Config struct {
	Redis    Redis
	DB       DB
	Rabbitmq RabbitMQ
	Forum    Forum
	PodID    string
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
	Support       Server
}

type Server struct {
	Host string
	Port string
}

type Redis struct {
	Host     string
	Port     string
	Password string
}

type RabbitMQ struct {
	Host     string
	Port     string
	Username string
	Password string
}
