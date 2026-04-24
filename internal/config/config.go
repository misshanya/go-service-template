package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HttpSrv  httpServer
	Postgres postgres
	S3       s3
	Kafka    kafka
	Valkey   valkey
	JWT      jwt
}

type httpServer struct {
	Addr string `env:"SERVER_ADDR" env-default:":8080"`
}

type postgres struct {
	URL      string `env:"POSTGRES_URL" env-required:"true"`
	MaxConns int32  `env:"POSTGRES_MAX_CONNS" env-default:"100"`
}

type s3 struct {
	Region     string `env:"S3_REGION" env-default:"us-east-1"`
	Endpoint   string `env:"S3_ENDPOINT"`
	AccessKey  string `env:"S3_ACCESS_KEY"`
	SecretKey  string `env:"S3_SECRET_KEY"`
	BucketName string `env:"S3_BUCKET_NAME"`
}

type kafka struct {
	Addr          string `env:"KAFKA_ADDR"`
	ReaderGroupID string `env:"KAFKA_READER_GROUP_ID"`
	Topic         string `env:"KAFKA_TOPIC"`
}

type valkey struct {
	Addr     string `env:"VALKEY_ADDR"`
	Password string `env:"VALKEY_PASSWORD"`
}

type jwt struct {
	Secret string `env:"JWT_SECRET" env-required:"true"`
}

func New() (*Config, error) {
	var cfg Config

	// Read .env file
	// If failed to read file, will try ReadEnv
	if err := cleanenv.ReadConfig(".env", &cfg); err == nil {
		return &cfg, nil
	}

	// Read env
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
