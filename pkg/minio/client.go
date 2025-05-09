package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log/slog"
	"time"
)

type Client struct {
	logger      *slog.Logger
	minioClient *minio.Client
}

// Создание самого клиента
func NewClient(endpoint, accessKeyID, secretAccessKey string, logger *slog.Logger) (*Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &Client{
		logger:      logger,
		minioClient: minioClient,
	}, nil
}

// Получение файла из бакета
func (c *Client) GetFile(ctx context.Context, bucketName, fileId string) (*minio.Object, error) {
	reqCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	obj, err := c.minioClient.GetObject(reqCtx, bucketName, fileId, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file %s:from minio bucket %s err: %w", fileId, bucketName, err)
	}
	return obj, nil
}

func (c *Client) UploadFile(ctx context.Context, fileId, fileName, bucketName string, fileSize int64, reader io.Reader) error {
	reqCtx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	exists, errBucketExists := c.minioClient.BucketExists(reqCtx, bucketName)
	if errBucketExists != nil || !exists {
		c.logger.Warn("no such bucket: " + bucketName + " creating new one...")
		err := c.minioClient.MakeBucket(reqCtx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			fmt.Println("createdF err")
			return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
		}
	}

	c.logger.Debug("put new object" + fileName + " in bucket " + bucketName)

	_, err := c.minioClient.PutObject(reqCtx, bucketName, fileId, reader, fileSize,
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name": fileName,
			},
		})
	if err != nil {
		return fmt.Errorf("failed to upload file. err: %w", err)
	}

	return nil
}

func (c *Client) DeleteFile(ctx context.Context, noteUUID, fileName string) error {
	err := c.minioClient.RemoveObject(ctx, noteUUID, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %w", fileName, err)
	}
	return nil
}
