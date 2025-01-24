package db

import (
	"log"

	"github.com/LikheKeto/Suraksheet/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient() *minio.Client {
	minioClient, err := minio.New(config.Envs.MinioURL, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Envs.MinioAccessKey, config.Envs.MinioSecretKey, ""),
		Secure: false, // TODO: make it https?
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Minio Successfully connected!")
	return minioClient
}
