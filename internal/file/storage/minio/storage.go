package minio

import (
	"bytes"
	"context"
	"fmt"
	"h-project/internal/file"
	"h-project/pkg/minio"
	"io"
	"log/slog"
)

type minioStorage struct {
	client *minio.Client
	logger *slog.Logger
}

func NewStorage(endpoint, accessKeyID, secretAccessKey string, logger *slog.Logger) (file.Storage, error) {
	client, err := minio.NewClient(endpoint, accessKeyID, secretAccessKey, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}
	return &minioStorage{client: client}, nil
}

func (m *minioStorage) GetFile(ctx context.Context, bucketName, fileID string) (*file.File, error) {
	obj, err := m.client.GetFile(ctx, bucketName, fileID)

	if err != nil {
		return nil, fmt.Errorf("failed to get the file: %w", err)
	}

	defer obj.Close()

	objectInfo, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get the file info: %w", err)
	}

	buffer := make([]byte, objectInfo.Size)
	_, err = obj.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to get objects: %w", err)
	}

	f := file.File{
		ID:    objectInfo.Key,
		Name:  objectInfo.UserMetadata["Name"],
		Size:  objectInfo.Size,
		Bytes: buffer,
	}

	return &f, nil
}

func (m *minioStorage) CreateFile(ctx context.Context, bucketName string, file *file.File) error {
	err := m.client.UploadFile(ctx, file.ID, file.Name, bucketName, file.Size, bytes.NewBuffer(file.Bytes))
	if err != nil {
		return err
	}
	return nil
}

func (m *minioStorage) DeleteFile(ctx context.Context, bucketName, fileID string) error {
	err := m.client.DeleteFile(ctx, bucketName, fileID)
	if err != nil {
		return err
	}
	return nil
}
