package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/ZephyrJung/QiaoYiCommunity/pkg/oss"
	"github.com/ZephyrJung/QiaoYiCommunity/utils"
)

type UploadService struct {
	storage oss.Storage
}

func NewUploadService(storage oss.Storage) *UploadService {
	return &UploadService{storage: storage}
}

const (
	maxFileSize = 5 * 1024 * 1024
	maxFiles    = 9
)

var allowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
}

func (s *UploadService) UploadImages(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	if len(files) == 0 {
		return nil, errors.New("no files uploaded")
	}
	if len(files) > maxFiles {
		return nil, fmt.Errorf("too many files, max %d", maxFiles)
	}

	urls := make([]string, 0, len(files))
	for _, file := range files {
		if file.Size > maxFileSize {
			return nil, fmt.Errorf("file %s exceeds max size 5MB", file.Filename)
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedExts[ext] {
			return nil, fmt.Errorf("file %s type not allowed", file.Filename)
		}

		f, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("open file %s failed: %w", file.Filename, err)
		}

		fileName := utils.GenerateIDString() + ext
		url, err := s.storage.Upload(ctx, fileName, f, file.Size)
		f.Close()
		if err != nil {
			return nil, fmt.Errorf("upload file %s failed: %w", file.Filename, err)
		}

		urls = append(urls, url)
	}

	return urls, nil
}

func (s *UploadService) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	urls, err := s.UploadImages(ctx, []*multipart.FileHeader{file})
	if err != nil {
		return "", err
	}
	return urls[0], nil
}
