package config

type Config struct {
	Port string
	DBPath string
	MigrationsPath string
}

func Load() *Config {
	return &Config {
		Port: ":4444",
		DBPath: "./data/realtime.db",
		MigrationsPath: "backend/db/migration",
	}
}