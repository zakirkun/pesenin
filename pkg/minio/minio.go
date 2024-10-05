package minio

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioContext struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

func (m *MinioContext) Open() (*minio.Client, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(m.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessKeyID, m.SecretAccessKey, ""),
		Secure: m.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func (m *MinioContext) Upload(bucketName, objectName, filePath string) {
	client, _ := m.Open()
	uploadFile(client, bucketName, objectName, filePath)
}

func (m *MinioContext) GetFile(bucketName, objectName string, expiry time.Duration) (string, error) {
	client, _ := m.Open()
	return getFileURL(client, bucketName, objectName, expiry)
}

func (m *MinioContext) Delete(bucketName, objectName string) {
	client, _ := m.Open()
	deleteFile(client, bucketName, objectName)
}

func uploadFile(minioClient *minio.Client, bucketName, objectName, filePath string) {
	// Make bucket if it doesn't exist.
	err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			log.Println(err)
		}
	}

	// Upload the file.
	info, err := minioClient.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		log.Println(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}

func getFileURL(minioClient *minio.Client, bucketName, objectName string, expiry time.Duration) (string, error) {
	// Generate a pre-signed URL for the object.
	reqParams := url.Values{} // Optionally, add custom request parameters here.
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func deleteFile(minioClient *minio.Client, bucketName, objectName string) {
	// Delete the object.
	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		log.Println(err)
	}

	log.Printf("Successfully deleted %s from bucket %s\n", objectName, bucketName)
}
