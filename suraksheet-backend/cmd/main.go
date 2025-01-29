package main

import (
	"database/sql"
	"log"

	"github.com/LikheKeto/Suraksheet/cmd/api"
	"github.com/LikheKeto/Suraksheet/config"
	"github.com/LikheKeto/Suraksheet/db"
	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func FatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// database
	database := db.NewSQLStorage(db.DBConfig{
		User:     config.Envs.DBUser,
		Host:     config.Envs.DBHost,
		Port:     config.Envs.DBPort,
		Password: config.Envs.DBPassword,
		DBname:   config.Envs.DBName,
	})
	initStorage(database)
	defer database.Close()

	// minio
	minioClient := db.NewMinioClient()

	// rabbitmq
	rmqConn, err := amqp.Dial(config.Envs.RabbitMQUrl)
	FatalIfErr(err)
	defer rmqConn.Close()
	log.Println("RabbitMQ successfully connected!")
	ch, err := rmqConn.Channel()
	FatalIfErr(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"extraction_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	FatalIfErr(err)

	// elasticsearch
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			config.Envs.ElasticsearchUrl,
		},
	})
	FatalIfErr(err)

	server := api.NewAPIServer(config.Envs.Port, database, minioClient, ch, q, esClient)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB Successfully connected!")
}
