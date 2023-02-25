package conf

type Config struct{}

func Init() (*Config, error) {
	cfg := &Config{}
	return cfg, nil
}
