package dto

type CreateAssessmentReq struct {
	OwnerType   string `json:"owner_type" validate:"required,oneof=course module lesson"`
	OwnerID     string `json:"owner_id" validate:"required,uuid4"`
	Type        string `json:"type" validate:"required,oneof=pre post quiz"`
	Title       string `json:"title" validate:"required"`
	PassScore   int    `json:"pass_score"`
	MaxAttempts *int   `json:"max_attempts"`
	TimeLimitS  *int   `json:"time_limit_s"`
}

type AddQuestionReq struct {
	Type        string  `json:"type" validate:"required,oneof=single_choice multiple_choice true_false short_text"`
	Stem        string  `json:"stem" validate:"required"`
	Explanation *string `json:"explanation"`
	Points      int     `json:"points"`
	Seq         int     `json:"seq"`
	Choices     []struct {
		Label     string `json:"label"`
		IsCorrect bool   `json:"is_correct"`
		Seq       int    `json:"seq"`
	} `json:"choices"`
}

type ListAssessmentsFilter struct {
	OwnerType string `json:"owner_type"`
	OwnerID   string `json:"owner_id"`
	Type      string `json:"type"`
}

type UpdateAssessmentReq struct {
	Title       *string `json:"title"`
	Type        *string `json:"type"` // "pretest" | "posttest" | "quiz" (ถ้าอยากล็อก type ก็ไม่ต้องรวม field นี้)
	PassScore   *int    `json:"pass_score"`
	TimeLimitS  *int    `json:"time_limit_s"`
	MaxAttempts *int    `json:"max_attempts"`
}

type UpdateQuestionReq struct {
	Stem   *string `json:"stem"`
	Points *int    `json:"points"`
	Seq    *int    `json:"seq"`
	// ถ้าอยากแก้ type ของคำถาม ให้เพิ่ม: Type *string `json:"type"`
}

type ChoiceUpsert struct {
	ID        *string `json:"id,omitempty"` // เวลาทำ replace-all: ใส่ ID ถ้าต้องแก้ของเดิม
	Label     string  `json:"label" validate:"required"`
	IsCorrect bool    `json:"is_correct"`
	Seq       int     `json:"seq"`
}

// Replace ทั้งชุดของ question (ถ้าไม่มี ID จะสร้างใหม่, ถ้ามี ID แล้วไม่ถูกส่งมา ถือว่า “ลบทิ้ง”)
type ReplaceChoicesReq struct {
	Choices []ChoiceUpsert `json:"choices" validate:"required"`
}
