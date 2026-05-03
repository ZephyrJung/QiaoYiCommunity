package oss

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type LocalStorage struct {
	basePath string
	domain   string
}

func NewLocalStorage(basePath, domain string) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
		domain:   domain,
	}
}

func (ls *LocalStorage) Upload(ctx context.Context, fileName string, reader io.Reader, size int64) (string, error) {
	now := time.Now()
	subDir := fmt.Sprintf("%04d%02d%02d", now.Year(), now.Month(), now.Day())
	dir := filepath.Join(ls.basePath, subDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create upload dir: %w", err)
	}

	fullPath := filepath.Join(dir, fileName)
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	relative := filepath.Join(subDir, fileName)
	if ls.domain != "" {
		return fmt.Sprintf("%s/%s", ls.domain, relative), nil
	}
	return "/" + relative, nil
}

func (ls *LocalStorage) Delete(ctx context.Context, fileName string) error {
	return os.Remove(filepath.Join(ls.basePath, fileName))
}
