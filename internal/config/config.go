// Package config encapsulates describing configuration for app
package config

// Сonfig contains config for app
type Сonfig struct {
	MongoUser     string `env:"MONGO_USER" envDefault:"root"`
	MongoPassword string `env:"MONGO_PASSWORD" envDefault:"password"`
	MongoHost     string `env:"MONGO_HOST" envDefault:"localhost"`
	MongoPort     int    `env:"MONGO_PORT" envDefault:"27017"`
	MongoDB       string `env:"MONGO_DB" envDefault:"test"`

	RedisPass string `env:"REDIS_PASS" envDefault:""`
	RedisHost string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort int    `env:"REDIS_PORT" envDefault:"6379"`

	RabbitUser string `env:"RABBIT_USER" envDefault:"guest"`
	RabbitPass string `env:"RABBIT_PASS" envDefault:"guest"`
	RabbitHost string `env:"RABBIT_HOST" envDefault:"localhost"`
	RabbitPort int    `env:"RABBIT_PORT" envDefault:"5672"`

	ConsumerNumber int `env:"CONSUMER_NUMBER" envDefault:"0"`
}
