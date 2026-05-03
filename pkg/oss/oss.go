package oss

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, fileName string, reader io.Reader, size int64) (url string, err error)
	Delete(ctx context.Context, fileName string) error
}
