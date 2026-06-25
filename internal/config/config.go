package config

type Config struct {
	App AppConfig `yaml:"app"`
	DB  DBConfig  `yaml:"db"`
}

type AppConfig struct {
	ServerPort string `yaml:"server_port"`
}

type DBConfig struct {
	DbAddr      string `yaml:"db_addr"`
	MaxOpenConn string
	MaxIdleConn string
	MaxIdleTime string
}
