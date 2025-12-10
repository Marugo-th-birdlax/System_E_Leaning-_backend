package service

import (
	"errors"
	"mime/multipart"

	"github.com/Marugo/birdlax/internal/modules/content/dto"
	"github.com/Marugo/birdlax/internal/modules/content/models"
	"github.com/google/uuid"
)

type svc struct {
	assetRepo  AssetRepo
	lessonRepo LessonRepo
	uploader   StorageUploader
}

func New(assetRepo AssetRepo, lessonRepo LessonRepo, uploader StorageUploader) Service {
	return &svc{assetRepo: assetRepo, lessonRepo: lessonRepo, uploader: uploader}
}

func (s *svc) UploadVideo(file *multipart.FileHeader) (*dto.UploadVideoResp, error) {
	if file == nil {
		return nil, errors.New("no file")
	}
	url, stored, size, mime, err := s.uploader.SaveVideo(file)
	if err != nil {
		return nil, err
	}

	a := &models.Asset{
		ID:        uuid.NewString(),
		Kind:      "video",
		Filename:  stored,
		MimeType:  mime,
		SizeBytes: size,
		Storage:   "local",
		URL:       url,
	}
	if err := s.assetRepo.CreateAsset(a); err != nil {
		return nil, err
	}
	return &dto.UploadVideoResp{
		AssetID:   a.ID,
		URL:       a.URL,
		Filename:  a.Filename,
		MimeType:  a.MimeType,
		SizeBytes: a.SizeBytes,
		Storage:   a.Storage,
	}, nil
}

func (s *svc) CreateLesson(req dto.CreateLessonReq) (*models.Lesson, error) {
	l := &models.Lesson{
		ID:           uuid.NewString(),
		ModuleID:     req.ModuleID,
		Title:        req.Title,
		ContentType:  req.ContentType,
		Seq:          req.Seq,
		AssetID:      req.AssetID,
		AssessmentID: req.AssessmentID,
		DurationS:    req.DurationS,
		IsMandatory:  true,
	}
	if err := s.lessonRepo.CreateLesson(l); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *svc) GetLesson(id string) (*models.Lesson, error) {
	return s.lessonRepo.GetByID(id)
}
func (s *svc) ListLessons(moduleID string, page, per int) ([]models.Lesson, int64, error) {
	return s.lessonRepo.List(moduleID, page, per)
}
func (s *svc) GetAsset(id string) (*models.Asset, error) {
	return s.assetRepo.GetByID(id)
}

func (s *svc) UpdateLesson(id string, req dto.UpdateLessonReq) (*models.Lesson, error) {
	l, err := s.lessonRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		l.Title = *req.Title
	}
	if req.ContentType != nil {
		l.ContentType = *req.ContentType
	}
	if req.Seq != nil {
		l.Seq = *req.Seq
	}
	if req.AssetID != nil {
		l.AssetID = req.AssetID // ถ้า field ใน model เป็น *string
	}
	if req.AssessmentID != nil {
		l.AssessmentID = req.AssessmentID // ถ้าเป็น *string
	}
	if req.DurationS != nil {
		l.DurationS = req.DurationS
	}
	if req.IsMandatory != nil {
		l.IsMandatory = *req.IsMandatory
	}

	if err := s.lessonRepo.UpdateLesson(l); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *svc) DeleteLesson(id string) error {
	return s.lessonRepo.DeleteLesson(id)
}
