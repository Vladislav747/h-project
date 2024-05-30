package file

import (
	"context"
	"log/slog"
)

type service struct {
	storage Storage
	logger  *slog.Logger
}

func NewService(storage Storage, logger *slog.Logger) (Service, error) {
	return &service{
		storage: storage,
		logger:  logger,
	}, nil
}

type Service interface {
	GetFile(ctx context.Context, bucketName, fileName string) (f *File, err error)
	Create(ctx context.Context, bucketName string, dto CreateFileDTO) error
	Delete(ctx context.Context, bucketName, fileName string) error
}

func (s *service) GetFile(ctx context.Context, noteUUID, fileName string) (f *File, err error) {
	f, err = s.storage.GetFile(ctx, noteUUID, fileName)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s *service) Create(ctx context.Context, bucketName string, dto CreateFileDTO) error {
	dto.NormalizeName()
	file, err := NewFile(dto)
	if err != nil {
		return err
	}
	err = s.storage.CreateFile(ctx, bucketName, file)
	//Передать по топику данные
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Delete(ctx context.Context, bucketName, fileName string) error {
	err := s.storage.DeleteFile(ctx, bucketName, fileName)
	if err != nil {
		return err
	}
	return nil
}
