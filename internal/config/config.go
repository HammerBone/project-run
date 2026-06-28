package config

type Config struct {
	App AppConfig
	DB  DBConfig
}

type AppConfig struct {
	ServerPort string `env:"APP_PORT"`
}

type DBConfig struct {
	DBAddr      string `env:"DB_ADDR"`
	MaxOpenConn string `env:"MAX_OPEN_CONN"`
	MaxIdleConn string `env:"MAX_IDLE_CONN"`
	MaxIdleTime string `env:"MAX_IDLE_TIME"`
}
