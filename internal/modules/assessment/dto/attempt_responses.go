package dto

type AttemptResp struct {
	ID           string   `json:"id"`
	AssessmentID string   `json:"assessment_id"`
	Status       string   `json:"status"`
	StartedAt    string   `json:"started_at"`
	SubmittedAt  *string  `json:"submitted_at,omitempty"`
	ScoreRaw     *int     `json:"score_raw,omitempty"`
	ScorePercent *float64 `json:"score_percent,omitempty"`
	IsPassed     *bool    `json:"is_passed,omitempty"`
}

type GradeResp struct {
	Attempt AttemptResp `json:"attempt"`
	// สรุป
	TotalQuestions int `json:"total_questions"`
	Correct        int `json:"correct"`
}
