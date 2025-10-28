package storage

import "mime/multipart"

type Uploader interface {
	SaveVideo(file *multipart.FileHeader) (url, storedName string, size int64, mime string, err error)
}
