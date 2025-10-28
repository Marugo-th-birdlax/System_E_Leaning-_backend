package dto

type EnrollCourseReq struct {
	CourseID string `json:"course_id" validate:"required"` // uuid
}

type StartLessonReq struct {
	// เผื่อใช้ภายหลังได้ (ตอนนี้ไม่บังคับ)
}

type TrackLessonReq struct {
	CurrentPosition int64 `json:"current_position" validate:"required,min=0"`
	MaxPosition     int64 `json:"max_position" validate:"required,min=1"`
}

type CompleteLessonReq struct {
	// optional: verify checksum ฯลฯ
}
