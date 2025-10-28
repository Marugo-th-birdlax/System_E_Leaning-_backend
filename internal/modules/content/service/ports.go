package service

import (
	"mime/multipart"

	"github.com/Marugo/birdlax/internal/modules/content/dto"
	"github.com/Marugo/birdlax/internal/modules/content/models"
)

type AssetRepo interface {
	CreateAsset(a *models.Asset) error
	GetByID(id string) (*models.Asset, error)
}

type LessonRepo interface {
	CreateLesson(l *models.Lesson) error
	GetByID(id string) (*models.Lesson, error)
	List(moduleID string, page, per int) ([]models.Lesson, int64, error)
}

type StorageUploader interface {
	SaveVideo(file *multipart.FileHeader) (url, storedName string, size int64, mime string, err error)
}

type Service interface {
	UploadVideo(file *multipart.FileHeader) (*dto.UploadVideoResp, error)
	CreateLesson(req dto.CreateLessonReq) (*models.Lesson, error)
	ListLessons(moduleID string, page, per int) ([]models.Lesson, int64, error)
	GetLesson(id string) (*models.Lesson, error)
	GetAsset(id string) (*models.Asset, error)
}
