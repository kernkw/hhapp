package config

// Config is the configuration struct
type Config struct {
	Addr string `envconfig:"ADDR" default:""`
	Port int    `envconfig:"PORT" default:"8080"`

	DBHost     string `envconfig:"DB_HOST" required:"true"`
	DBPort     int    `envconfig:"DB_PORT" default:"3306"`
	DBName     string `envconfig:"DB_NAME" default:"pps"`
	DBUser     string `envconfig:"DB_USER" required:"true"`
	DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
}
