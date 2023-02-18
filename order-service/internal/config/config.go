package config

import (
	"order-service/pkg/broker"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App       AppConfig      `yaml:"app" json:"app"`
	HTTP      HTTPConfig     `yaml:"http" json:"HTTP"`
	Log       LogConfig      `yaml:"logger" json:"log"`
	Databases DatabaseConfig `yaml:"databases" json:"databases"`
	Timeout   TimeoutConfig  `yaml:"timeout" json:"timeout"`
	Cache     RedisConfig    `yaml:"redis" json:"redis"`
	Topics    Topics         `yaml:"topics" json:"topics"`
	Kafka     KafkaConfig    `yaml:"kafka" json:"kafka" env-prefix:"KAFKA_"`
}

type AppConfig struct {
	Name    string `env:"APP_NAME" env-required:"true" yaml:"name"`
	Version string `env:"APP_VERSION" env-required:"true" yaml:"version"`
}

type HTTPConfig struct {
	Port string `env:"APP_PORT" env-required:"true" yaml:"port"`
}

type LogConfig struct {
	Level string `env:"APP_LOGLEVEL" env-required:"true" yaml:"log_level"`
}

type TimeoutConfig struct {
	Duration int `env:"APP_TIMEOUT" env-required:"true" yaml:"duration"`
}

type ConnConfig struct {
	// Network the network type, either tcp or unix. Default is tcp.
	Network string `yaml:"network"`

	// Addr the database server host.
	Host string `yaml:"host" env:"HOST" env-required:"true"`

	// Port the database server port.
	Port string `yaml:"port" env:"PORT" env-required:"true"`

	// Username to authenticate the current connection.
	Username string `yaml:"username" env:"USER"`

	// Password must match the password specified in the
	// requirement pass server configuration option.
	Password string `yaml:"password" env:"PASSWORD"`

	Name string `yaml:"dbname" env:"NAME"`
}

type SQLConfig struct {
	ConnConfig `yaml:"conn_config"`

	// MaxOpenConns the maximum number of open connections to the database.
	//
	// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
	// MaxIdleConns, then MaxIdleConns will be reduced to match the new
	// MaxOpenConns limit.
	//
	// If n <= 0, then there is no limit on the number of open connections.
	// The default is 0 (unlimited).
	MaxOpenConns int `yaml:"max_open_conns"`

	// MaxIdleConns the maximum number of connections in the idle
	// connection pool.
	//
	// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns,
	// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit.
	//
	// If n <= 0, no idle connections are retained.
	//
	// The default max idle connections is currently 2. This may change in
	// a future.
	MaxIdleConns int `yaml:"max_idle_conns"`

	// ConnMaxIdleTime the maximum amount of time a connection may be idle.
	//
	// Expired connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are not closed due to a connection's idle time.
	// The default is 0 (unlimited).
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`

	// ConnMaxLifetime the maximum amount of time a connection may be reused.
	//
	// Expired connections may be closed lazily before reuse.
	//
	// If d <= 0, connections are not closed due to a connection's age.
	// The default is 0 (unlimited).
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

type DatabaseConfig struct {
	Postgres SQLConfig `yaml:"postgres" env-prefix:"POSTGRES_"`
}

type KafkaConfig struct {
	// BrokerAddresses initial list of brokers as a list of broker host in host:port format.
	BrokerAddresses []string `yaml:"brokers" env-required:"true"`
	SASL            struct {
		Username string `yaml:"username" env:"USER" env-required:"true"`
		Password string `yaml:"password" env:"PASSWORD"`
	} `yaml:"sasl"`
	SSL struct {
		CALocation string `yaml:"ca-location"`
	} `yaml:"ssl"`
}

func (c *KafkaConfig) ToSDKFormat() *broker.NetworkConfig {
	cfg := &broker.NetworkConfig{
		BrokerAddresses: c.BrokerAddresses,
		SASL: &broker.SASL{
			Username: c.SASL.Username,
			Password: c.SASL.Password,
		},
	}

	if c.SSL.CALocation != "" {
		cfg.SSL = &broker.SSL{
			CALocation: c.SSL.CALocation,
		}
	}

	return cfg
}

type RedisConfig struct {
	ConnConfig `yaml:"conn_config" env-prefix:"REDIS_"`
	DB         int `yaml:"db_number"`
}

type Topics struct {
	BillingTopic string `yaml:"billing_topic" env-default:"hire_requests"`
	StockTopic   string `yaml:"stock_topic" env-default:"hire_responses"`
	OrderTopic   string `yaml:"order_topic" env-default:"tasks_run"`
}

func ReadConfig(filePath string) (Config, error) {
	cfg := Config{}

	if err := cleanenv.ReadConfig(filePath, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
