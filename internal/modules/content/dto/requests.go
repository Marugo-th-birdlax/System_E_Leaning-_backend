package dto

type CreateLessonReq struct {
	ModuleID     string  `json:"module_id" validate:"required,uuid4"`
	Title        string  `json:"title" validate:"required"`
	ContentType  string  `json:"content_type" validate:"required,oneof=slide video document quiz"`
	Seq          int     `json:"seq" validate:"required,min=1"`
	AssetID      *string `json:"asset_id"`      // เมื่อเป็น video/slide/doc
	AssessmentID *string `json:"assessment_id"` // เมื่อเป็น quiz
	DurationS    *int64  `json:"duration_s"`
}

type UploadVideoResp struct {
	AssetID   string `json:"asset_id"`
	URL       string `json:"url"`
	Filename  string `json:"filename"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
	Storage   string `json:"storage"`
}
