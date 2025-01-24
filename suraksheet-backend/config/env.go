package config

import (
	"fmt"
	"os"
)

type Config struct {
	PublicHost string
	Port       string
	JWTSecret  string

	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string

	MinioURL        string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string

	RabbitMQUrl string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		PublicHost: getEnv("SERVER_PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("SERVER_PORT", ":8080"),
		JWTSecret:  getEnv("SERVER_JWT_SECRET", ""),
		DBUser:     getEnv("MYSQL_USER", "root"),
		DBPassword: getEnv("MYSQL_PASSWORD", "mypassword"),
		DBAddress: fmt.Sprintf("%s:%s", getEnv("MYSQL_HOST", "127.0.0.1"),
			getEnv("MYSQL_PORT", "3306")),
		DBName:          getEnv("MYSQL_DATABASE", "suraksheet"),
		MinioURL:        getEnv("MINIO_ENDPOINT", "127.0.0.1:9000"),
		MinioAccessKey:  getEnv("MINIO_ACCESS_KEY", ""),
		MinioSecretKey:  getEnv("MINIO_SECRET_KEY", ""),
		MinioBucketName: getEnv("MINIO_BUCKET_NAME", "suraksheet"),
		RabbitMQUrl:     getEnv("RABBITMQ_URL", "localhost"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
