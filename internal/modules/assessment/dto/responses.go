package dto

type PageMeta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
}

// สำหรับตอบกลับหลังสร้าง Assessment
type AssessmentResp struct {
	ID          string `json:"id"`
	OwnerType   string `json:"owner_type"`
	OwnerID     string `json:"owner_id"`
	Type        string `json:"type"` // pre|post|quiz
	Title       string `json:"title"`
	PassScore   int    `json:"pass_score"`
	MaxAttempts *int   `json:"max_attempts,omitempty"`
	TimeLimitS  *int   `json:"time_limit_s,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}
type AssessmentItem struct {
	ID          string `json:"id"`
	OwnerType   string `json:"owner_type"`
	OwnerID     string `json:"owner_id"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	PassScore   int    `json:"pass_score"`
	TimeLimitS  *int   `json:"time_limit_s,omitempty"`
	MaxAttempts *int   `json:"max_attempts,omitempty"`
}

// สำหรับตอบกลับหลังเพิ่มคำถาม
type QuestionResp struct {
	ID           string       `json:"id"`
	AssessmentID string       `json:"assessment_id,omitempty"` // <- เพิ่มบรรทัดนี้ถ้าต้องการ
	Type         string       `json:"type"`
	Stem         string       `json:"stem"`
	Points       int          `json:"points"`
	Seq          int          `json:"seq"`
	Choices      []ChoiceResp `json:"choices,omitempty"`
}

type PagedAssessments struct {
	Data []AssessmentItem `json:"data"`
	Meta PageMeta         `json:"meta"`
}

type ChoiceResp struct {
	ID        string `json:"id"`
	Label     string `json:"label"`
	IsCorrect bool   `json:"is_correct"`
	Seq       int    `json:"seq"`
}

type AssessmentDetailResp struct {
	Assessment AssessmentItem `json:"assessment"`
	Questions  []QuestionResp `json:"questions"`
}
