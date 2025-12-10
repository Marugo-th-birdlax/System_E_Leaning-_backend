package service

import (
	"errors"
	"math"
	"time"

	lmodels "github.com/Marugo/birdlax/internal/modules/learning/models"
	learningrepo "github.com/Marugo/birdlax/internal/modules/learning/repo"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MetricsService interface {
	OnSubmitAttempt(userID, courseID string, score float64, passed bool, timeSec int64) error
	OnCourseCompleted(userID, courseID string, timeSec int64) error
	GetLearningMetric(userID, courseID string) (*lmodels.LearningMetric, error)
	GetCourseOutcome(courseID string) (*lmodels.CourseOutcome, error)
}

type metricsSvc struct {
	repo *learningrepo.MetricsRepo
}

func NewMetricsService(r *learningrepo.MetricsRepo) MetricsService {
	return &metricsSvc{repo: r}
}

func (s *metricsSvc) GetLearningMetric(userID, courseID string) (*lmodels.LearningMetric, error) {
	return s.repo.GetLearningMetric(userID, courseID)
}
func (s *metricsSvc) GetCourseOutcome(courseID string) (*lmodels.CourseOutcome, error) {
	return s.repo.GetCourseOutcome(courseID)
}

// OnSubmitAttempt updates per-user metric and lightly nudges course outcome (approx)
func (s *metricsSvc) OnSubmitAttempt(userID, courseID string, score float64, passed bool, timeSec int64) error {
	now := time.Now()
	m, err := s.repo.GetLearningMetric(userID, courseID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// real error -> return it so caller can know
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		m = nil
	}

	if m == nil {
		m = &lmodels.LearningMetric{
			ID:               uuid.NewString(),
			UserID:           userID,
			CourseID:         courseID,
			AvgScore:         score,
			LastScore:        &score,
			AttemptsCount:    1,
			PassCount:        0,
			TotalTimeSeconds: timeSec,
			CompletionStatus: "in_progress",
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		if passed {
			m.PassCount = 1
		}
		return s.repo.UpsertLearningMetric(m)
	}

	// incremental average
	prevAttempts := m.AttemptsCount
	newAttempts := prevAttempts + 1
	m.AvgScore = (m.AvgScore*float64(prevAttempts) + score) / float64(newAttempts)
	m.AttemptsCount = newAttempts
	m.LastScore = &score
	if passed {
		m.PassCount = m.PassCount + 1
	}
	m.TotalTimeSeconds = m.TotalTimeSeconds + timeSec
	if m.CompletionStatus == "" {
		m.CompletionStatus = "in_progress"
	}
	m.UpdatedAt = now

	if err := s.repo.UpsertLearningMetric(m); err != nil {
		return err
	}

	// light update for course outcome: keep simple incremental values
	co, err := s.repo.GetCourseOutcome(courseID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// ถ้ามี error จริงๆ ให้คืน error (หรืออย่างน้อย log แล้ว continue)
		// ผมเลือก log + return nil เพื่อไม่ทำให้ submit fail — แต่ต้องมี log
		// log.Printf("GetCourseOutcome error: %v", err)
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		co = nil
	}
	if co == nil {
		co = &lmodels.CourseOutcome{
			CourseID:          courseID,
			TotalEnrollments:  0,
			TotalCompleted:    0,
			AvgScore:          m.AvgScore,
			PassRate:          0,
			MedianTimeSeconds: 0,
			UpdatedAt:         now,
		}
		if err := s.repo.UpsertCourseOutcome(co); err != nil {
			return nil
		}
		return nil
	}
	// approximate avg across users: simple smoothing (not exact). If you want precise, batch recompute.
	co.AvgScore = (co.AvgScore + m.AvgScore) / 2.0
	// pass rate approximation
	co.PassRate = math.Min(100, (co.PassRate+mPct(m))/2.0)
	co.UpdatedAt = now
	_ = s.repo.UpsertCourseOutcome(co)
	return nil
}

func mPct(m *lmodels.LearningMetric) float64 {
	if m.AttemptsCount == 0 {
		return 0
	}
	return (float64(m.PassCount) / float64(m.AttemptsCount)) * 100.0
}

// OnCourseCompleted: mark user metric completed and bump course_outcome.total_completed
func (s *metricsSvc) OnCourseCompleted(userID, courseID string, timeSec int64) error {
	now := time.Now()
	m, _ := s.repo.GetLearningMetric(userID, courseID)
	if m == nil {
		m = &lmodels.LearningMetric{
			ID:               uuid.NewString(),
			UserID:           userID,
			CourseID:         courseID,
			AvgScore:         0,
			LastScore:        nil,
			AttemptsCount:    0,
			PassCount:        0,
			TotalTimeSeconds: timeSec,
			CompletionStatus: "completed",
			CreatedAt:        now,
			UpdatedAt:        now,
		}
	} else {
		m.CompletionStatus = "completed"
		m.TotalTimeSeconds = m.TotalTimeSeconds + timeSec
		m.UpdatedAt = now
	}
	if err := s.repo.UpsertLearningMetric(m); err != nil {
		return err
	}

	co, _ := s.repo.GetCourseOutcome(courseID)
	if co == nil {
		co = &lmodels.CourseOutcome{
			CourseID:          courseID,
			TotalEnrollments:  0,
			TotalCompleted:    1,
			AvgScore:          m.AvgScore,
			PassRate:          mPct(m),
			MedianTimeSeconds: m.TotalTimeSeconds,
			UpdatedAt:         now,
		}
		return s.repo.UpsertCourseOutcome(co)
	}
	co.TotalCompleted = co.TotalCompleted + 1
	// naive avg update:
	if co.AvgScore == 0 {
		co.AvgScore = m.AvgScore
	} else {
		co.AvgScore = (co.AvgScore + m.AvgScore) / 2.0
	}
	co.PassRate = math.Min(100, (co.PassRate+mPct(m))/2.0)
	co.UpdatedAt = now
	return s.repo.UpsertCourseOutcome(co)
}
