package config

type Config struct {
	ConnUrl string
}

func New(connUrl string) Config {
	return Config{
		ConnUrl: connUrl,
	}
}