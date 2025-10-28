package dto

type EnrollmentResp struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	CourseID        string  `json:"course_id"`
	Status          string  `json:"status"`
	ProgressPercent float64 `json:"progress_percent"`
}

type LessonProgressResp struct {
	UserID          string  `json:"user_id"`
	LessonID        string  `json:"lesson_id"`
	ProgressPercent float64 `json:"progress_percent"`
	CurrentPosition int64   `json:"current_position"`
	MaxPosition     int64   `json:"max_position"`
	IsUnlocked      bool    `json:"is_unlocked"`
}
