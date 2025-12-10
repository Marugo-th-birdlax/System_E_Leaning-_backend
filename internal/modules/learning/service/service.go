package service

import (
	"errors"
	"time"

	"github.com/Marugo/birdlax/internal/modules/learning/dto"
	"github.com/Marugo/birdlax/internal/modules/learning/models"
	"github.com/google/uuid"
)

type svc struct {
	repo    Repo
	metrics MetricsService
}

func New(r Repo, ms MetricsService) Service { return &svc{repo: r, metrics: ms} }

func (s *svc) EnrollCourse(userID, courseID string) (*models.Enrollment, error) {
	now := time.Now()
	e := &models.Enrollment{
		ID:             uuid.NewString(),
		UserID:         userID,
		CourseID:       courseID,
		Status:         models.StatusInProgress,
		StartedAt:      &now,
		LastAccessedAt: &now,
	}

	return e, s.repo.UpsertEnrollment(e)
}

func (s *svc) GetEnrollment(userID, courseID string) (*models.Enrollment, error) {
	return s.repo.GetEnrollment(userID, courseID)
}

// ========== LESSON FLOW ==========

func (s *svc) StartLesson(userID, lessonID string, _ dto.StartLessonReq) (*models.UserLessonProgress, error) {
	lesson, err := s.repo.GetLesson(lessonID)
	if err != nil {
		return nil, err
	}

	// Anti-skip: ถ้า seq > 1 ต้องผ่าน seq-1 ในโมดูลเดียวกัน
	if lesson.Seq > 1 {
		prev, err := s.repo.GetPrevLessonSameModule(lesson.ModuleID, lesson.Seq)
		if err != nil {
			return nil, errors.New("cannot start: missing previous lesson")
		}
		if prev != nil {
			prevProg, err := s.repo.GetLessonProgress(userID, prev.ID)
			if err != nil || prevProg == nil || prevProg.CompletedAt == nil {
				return nil, errors.New("cannot start: previous lesson not completed")
			}
		}
	}

	// Progress record
	p, err := s.repo.GetLessonProgress(userID, lessonID)
	now := time.Now()
	if err != nil {
		// สร้างใหม่
		p = &models.UserLessonProgress{
			ID:              uuid.NewString(),
			UserID:          userID,
			LessonID:        lessonID,
			IsUnlocked:      true,
			StartedAt:       &now,
			ProgressPercent: 0,
			CurrentPosition: 0,
			MaxPosition:     0,
		}
		if err := s.repo.CreateLessonProgress(p); err != nil {
			return nil, err
		}
		return p, nil
	}

	// มีอยู่แล้ว → mark unlock + started (ครั้งแรก)
	if !p.IsUnlocked {
		p.IsUnlocked = true
	}
	if p.StartedAt == nil {
		p.StartedAt = &now
	}
	if err := s.repo.UpdateLessonProgress(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *svc) TrackLesson(userID, lessonID string, req dto.TrackLessonReq) (*models.UserLessonProgress, error) {
	p, err := s.repo.GetLessonProgress(userID, lessonID)
	if err != nil {
		return nil, errors.New("lesson not started")
	}

	// Resume + Percent
	p.CurrentPosition = req.CurrentPosition
	if req.MaxPosition > p.MaxPosition {
		p.MaxPosition = req.MaxPosition
	}
	if p.MaxPosition > 0 && p.CurrentPosition >= 0 {
		ratio := float64(p.CurrentPosition) / float64(p.MaxPosition)
		if ratio < 0 {
			ratio = 0
		}
		if ratio > 1 {
			ratio = 1
		}
		p.ProgressPercent = ratio * 100
	}

	if err := s.repo.UpdateLessonProgress(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *svc) CompleteLesson(userID, lessonID string, _ dto.CompleteLessonReq) (*models.UserLessonProgress, error) {
	p, err := s.repo.GetLessonProgress(userID, lessonID)
	if err != nil {
		return nil, errors.New("lesson not started")
	}

	now := time.Now()
	p.ProgressPercent = 100
	p.CompletedAt = &now
	if err := s.repo.UpdateLessonProgress(p); err != nil {
		return nil, err
	}

	// note: handler already calls UpdateEnrollmentPercent afterwards in your code.
	// alternative: you could call UpdateEnrollmentPercent here if you have courseID available.

	return p, nil
}

func (s *svc) UpdateEnrollmentPercent(userID, courseID string) error {
	total, err := s.repo.CountMandatoryLessonsOfCourse(courseID)
	if err != nil {
		return err
	}
	if total == 0 {
		return nil
	}
	done, err := s.repo.CountCompletedMandatoryLessons(userID, courseID)
	if err != nil {
		return err
	}
	percent := float64(done) / float64(total) * 100.0

	e, err := s.repo.GetEnrollment(userID, courseID)
	if err != nil {
		return err
	}
	now := time.Now()
	e.ProgressPercent = percent
	e.LastAccessedAt = &now

	if percent >= 100 {
		// ถ้ายังไม่เคย completed/passed/failed → mark เป็น completed
		if e.Status == "" ||
			e.Status == models.StatusEnrolled ||
			e.Status == models.StatusInProgress {
			e.Status = models.StatusCompleted
		}
		if e.CompletedAt == nil {
			e.CompletedAt = &now
		}

		// call metrics service if available (guard against nil)
		if s.metrics != nil {
			// ให้เวลาเป็น 0 ถ้าไม่มีการเก็บเวลาจริง; ถ้าต้องการ คำนวณจาก lesson progress timestamps
			var totalSec int64 = 0
			_ = s.metrics.OnCourseCompleted(userID, courseID, totalSec)
		}
	}

	return s.repo.UpsertEnrollment(e)
}
