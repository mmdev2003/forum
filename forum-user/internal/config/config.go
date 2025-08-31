package config

func New() *Config {
	return &Config{
		DB: DB{
			Username: "postgres",
			Dbname:   "postgres",
			Password: "postgres",
			Host:     "forum-user-db",
			Port:     "5432",
		},
		Meilisearch: Meilisearch{
			"forum-meilisearch",
			"7700",
			"qwertyuiop",
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
			User: Server{
				Host: "forum-user",
				Port: "8006",
			},
			Notification: Server{
				Host: "forum-notification",
				Port: "8007",
			},
		},
	}
}

type Config struct {
	DB          DB
	Meilisearch Meilisearch
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
	Authorization Server
	Notification  Server
	User          Server
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

type Server struct {
	Host string
	Port string
}
