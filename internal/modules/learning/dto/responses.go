package dto

type EnrollmentResp struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	CourseID        string  `json:"course_id"`
	Status          string  `json:"status"`
	ProgressPercent float64 `json:"progress_percent"`
}

type LessonProgressResp struct {
	LessonID        string  `json:"lesson_id"`
	Title           string  `json:"title"`
	ModuleID        string  `json:"module_id"`
	Seq             int     `json:"seq"`
	ContentType     string  `json:"content_type"`
	ProgressPercent float64 `json:"progress_percent"`
	CurrentPosition int64   `json:"current_position"`
	MaxPosition     int64   `json:"max_position"`
	IsUnlocked      bool    `json:"is_unlocked"`
	StartedAt       *string `json:"started_at,omitempty"`
	CompletedAt     *string `json:"completed_at,omitempty"`
}

type CourseProgressResp struct {
	Enrollment *EnrollmentResp      `json:"enrollment,omitempty"`
	Lessons    []LessonProgressResp `json:"lessons"`
}
