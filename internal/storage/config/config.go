package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// export $(grep -v '^#' .env | xargs) to read env from dev file
type Config struct {
	EnvName            string `env:"ENV_NAME" env-required:"true" env-description:"env name"`
	LogLevel           string `env:"LOG_LEVEL" env-default:"info" env-description:"logging level, debug, info, etc."`
	ApiVersion         string `env:"API_VERSION" env-required:"true" env-description:"API version"`
	Server             HttpServer
	PrivateHttpServer  HttpServer
	Database           Database
	Cache              Cache
	StandardHttpClient StandardHttpClient
	MinioStorage       MinioStorage
}

type StandardHttpClient struct {
	Timeout time.Duration `env:"STANDARD_HTTP_CLIENT_TIMEOUT" env-default:"3s" env-description:"default timeout"`
}

type HttpServer struct {
	Host           string        `env:"SERVER_HOST" env-default:"" env-description:"tcp network address"`
	Port           int           `env:"SERVER_PORT" env-default:"80" env-description:"server port"`
	DebugPort      int           `env:"SERVER_DEBUG_PORT" env-default:"8081" env-description:"metrics scrapping port"`
	CloseTimeout   time.Duration `env:"SERVER_GRACEFUL_SHUTDOWN_TIMEOUT" env-default:"10s" env-description:"timeout for http graceful shutdown"`
	SwaggerEnabled bool          `env:"SERVER_SWAGGER_ENABLED" env-default:"false" env-description:"writes error-level logs to sentry"`
}

type Database struct {
	Net                string        `env:"MYSQL_NET" env-default:"tcp"`
	Server             string        `env:"MYSQL_SERVER" env-required:"true" env-description:"host and port, e.g. localhost:3306"`
	DBName             string        `env:"MYSQL_DBNAME" env-required:"true" env-description:"database name, e.g. dostavkee"`
	User               string        `env:"MYSQL_USER" env-default:"root" env-description:"database username, e.g. dostavkee"`
	Password           string        `env:"MYSQL_PASSWORD" env-default:"root" env-description:"database password, e.g. notasecret"`
	TimeZone           string        `env:"MYSQL_TIMEZONE" env-default:"UTC" env-description:"database timezone"`
	Timeout            time.Duration `env:"MYSQL_RECONNECT_TIMEOUT" env-default:"2s" env-description:"how much trying to connect to db"`
	MaxIdleConnections int           `env:"MYSQL_MAX_IDLE_CONNS" env-default:"40" env-description:"how much idle tcp conns are possible"`
	MaxOpenConnections int           `env:"MYSQL_MAX_OPEN_CONNS" env-default:"40" env-description:"max tcp conns to the database"`
}

type Cache struct {
	Type  string `env:"REDIS_TYPE" env-required:"true" env-description:"specifies provider, one of redis/redisCluster"`
	Redis struct {
		Address  string `env:"REDIS_ADDR" env-default:"" env-description:"redis host:port single instance"`
		Password string `env:"REDIS_PASSWORD" env-default:"" env-description:"redis password if exists"`
		PoolSize int    `env:"REDIS_POOL_SIZE" env-default:"70" env-description:"max tcp connections pool size"`
	}
	RedisCluster struct {
		Addresses []string `env:"REDIS_CLUSTER_ADDRS" env-default:"" env-description:"redis cluster nodes: ['172.27.29.90:7000','172.27.29.91:7001'', '172.27.29.92:7002'']"`
		Password  string   `env:"REDIS_PASSWORD" env-default:"" env-description:"redis password if exists"`
		PoolSize  int      `env:"REDIS_POOL_SIZE" env-default:"70" env-description:"max tcp connections pool size"`
	}
}

type MinioStorage struct {
	Host      string `env:"MINIO_STORAGE_HOST" env-required:"true" env-description:"Minio storage host"`
	AccessKey string `env:"MINIO_STORAGE_ACCESS_KEY" env-required:"true" env-description:"Minio storage access key"`
	SecretKey string `env:"MINIO_STORAGE_SECRET_KEY" env-required:"true" env-description:"Minio storage secret key"`
	Bucket    string `env:"MINIO_STORAGE_BUCKET" env-required:"true" env-description:"Minio storage bucket"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
