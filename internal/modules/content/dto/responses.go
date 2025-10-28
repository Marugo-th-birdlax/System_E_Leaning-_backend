package dto

type LessonResp struct {
	ID           string  `json:"id"`
	ModuleID     string  `json:"module_id"`
	Title        string  `json:"title"`
	ContentType  string  `json:"content_type"`
	Seq          int     `json:"seq"`
	AssetID      *string `json:"asset_id"`
	AssessmentID *string `json:"assessment_id"`
	DurationS    *int64  `json:"duration_s"`
}
