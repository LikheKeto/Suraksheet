package utils

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/LikheKeto/Suraksheet/config"
	"github.com/minio/minio-go/v7"
)

// List of allowed content types
var allowedContentTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"application/pdf": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

func GetObject(ctx context.Context, minioClient *minio.Client, objectName string) (*minio.Object, error) {
	obj, err := minioClient.GetObject(ctx, config.Envs.MinioBucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, err
}

func DeleteObject(ctx context.Context, minioClient *minio.Client, object string) error {
	return minioClient.RemoveObject(ctx, config.Envs.MinioBucketName, object, minio.RemoveObjectOptions{})
}

func DeleteDir(ctx context.Context, minioClient *minio.Client, dir string) error {
	toDeleteChan := make(chan minio.ObjectInfo)
	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(toDeleteChan)
		for object := range minioClient.ListObjects(context.Background(), config.Envs.MinioBucketName, minio.ListObjectsOptions{
			Prefix:    dir,
			Recursive: true,
		}) {
			if object.Err != nil {
				log.Fatalln(object.Err)
			}
			toDeleteChan <- object
		}
	}()
	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}
	for rErr := range minioClient.RemoveObjects(ctx, config.Envs.MinioBucketName, toDeleteChan, opts) {
		fmt.Println("Error detected during deletion: ", rErr)
	}
	return nil
}

func RenameObject(ctx context.Context, minioClient *minio.Client, old, new string) error {
	_, err := minioClient.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: config.Envs.MinioBucketName,
		Object: new,
	}, minio.CopySrcOptions{
		Bucket: config.Envs.MinioBucketName,
		Object: old,
	})
	return err
}

func UploadToMinio(ctx context.Context, minioClient *minio.Client,
	file multipart.File, fileHeader *multipart.FileHeader,
	referenceName string) error {
	defer file.Close()
	// Define the bucket name and object name
	bucketName := config.Envs.MinioBucketName
	objectName := referenceName
	contentType := fileHeader.Header.Get("Content-Type")

	// Check if content type is allowed
	if !allowedContentTypes[contentType] {
		return fmt.Errorf("file type not allowed: %s", contentType)
	}

	// Upload the file to MinIO
	_, err := minioClient.PutObject(ctx, bucketName,
		objectName, file, fileHeader.Size,
		minio.PutObjectOptions{ContentType: contentType})
	return err
}
