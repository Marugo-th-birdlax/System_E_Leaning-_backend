package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type LocalFS struct {
	BaseDir string // e.g. "./uploads/videos"
	BaseURL string // e.g. "/static/videos" (เสิร์ฟผ่าน Nginx/Fiber Static)
}

func (l *LocalFS) SaveVideo(file *multipart.FileHeader) (string, string, int64, string, error) {
	fn := uuid.NewString() + filepath.Ext(file.Filename)
	dstPath := filepath.Join(l.BaseDir, fn)

	src, err := file.Open()
	if err != nil {
		return "", "", 0, "", err
	}
	defer src.Close()

	if err := os.MkdirAll(l.BaseDir, 0755); err != nil {
		return "", "", 0, "", err
	}
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", "", 0, "", err
	}
	defer dst.Close()

	n, err := io.Copy(dst, src)
	if err != nil {
		return "", "", 0, "", err
	}

	url := fmt.Sprintf("%s/%s?ts=%d", l.BaseURL, fn, time.Now().Unix())
	return url, fn, n, file.Header.Get("Content-Type"), nil
}
