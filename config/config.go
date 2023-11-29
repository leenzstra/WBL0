package config

type Config struct {
	PgUrl string
	NatsUrl string
}

func New(pgUrl string, natsUrl string) Config {
	return Config{
		PgUrl: pgUrl,
		NatsUrl: natsUrl,
	}
}