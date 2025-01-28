package config

import (
	"os"
	"strconv"
)

type Config struct {
	PublicHost string
	Port       string
	JWTSecret  string

	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
	DBPort     int

	MinioURL        string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucketName string

	RabbitMQUrl string
}

var Envs = initConfig()

func initConfig() Config {
	port, err := strconv.Atoi(getEnv("POSTGRES_PORT", "5432"))
	if err != nil {
		panic(err)
	}
	return Config{
		PublicHost:      getEnv("SERVER_PUBLIC_HOST", "http://localhost"),
		Port:            getEnv("SERVER_PORT", ":8080"),
		JWTSecret:       getEnv("SERVER_JWT_SECRET", ""),
		DBUser:          getEnv("POSTGRES_USER", "root"),
		DBPassword:      getEnv("POSTGRES_PASSWORD", "mypassword"),
		DBHost:          getEnv("POSTGRES_HOST", "127.0.0.1"),
		DBPort:          port,
		DBName:          getEnv("POSTGRES_DATABASE", "suraksheet"),
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
