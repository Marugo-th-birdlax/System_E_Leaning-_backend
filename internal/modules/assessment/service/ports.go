package service

import (
	"time"

	"github.com/Marugo/birdlax/internal/modules/assessment/dto"
	"github.com/Marugo/birdlax/internal/modules/assessment/models"
	assrepo "github.com/Marugo/birdlax/internal/modules/assessment/repo"
)

// Repo คือพอร์ตฝั่ง persistence (GORM)
type Repo interface {
	CreateAssessment(a *models.Assessment) error
	AddQuestion(q *models.Question, choices []models.Choice) error
	ListAssessments(ownerType, ownerID, atype string, page, per int) ([]models.Assessment, int64, error)
	GetAssessmentWithItems(id string) (*models.Assessment, []models.Question, map[string][]models.Choice, error)

	UpdateAssessment(a *models.Assessment) error
	DeleteAssessment(id string) error

	UpdateQuestion(q *models.Question) error
	DeleteQuestion(id string) error
	GetQuestionByID(id string) (*models.Question, error)
	GetAssessmentByID(id string) (*models.Assessment, error)
}

// Service คือ business use cases
type Service interface {
	Create(req dto.CreateAssessmentReq) (*models.Assessment, error)
	AddQuestion(assessID string, req dto.AddQuestionReq) (*models.Question, error)

	List(filter dto.ListAssessmentsFilter, page, per int) ([]dto.AssessmentItem, int64, error)
	GetDetail(id string) (*dto.AssessmentDetailResp, error)
	UpdateAssessment(id string, req dto.UpdateAssessmentReq) (*models.Assessment, error)
	DeleteAssessment(id string) error

	UpdateQuestion(id string, req dto.UpdateQuestionReq) (*models.Question, error)
	DeleteQuestion(id string) error
}

type AttemptRepo interface {
	CountUserAttempts(assessmentID, userID string) (int64, error)
	CreateAttempt(a *models.Attempt) error
	GetAttempt(id, userID string) (*models.Attempt, error)
	UpdateAttempt(a *models.Attempt) error

	UpsertAnswer(ans *models.Answer) error
	ListAnswers(attemptID string) ([]models.Answer, error)

	ListQuestionsWithCorrectChoices(assessmentID string) ([]assrepo.QuestionWithChoices, error)
	Now() (t time.Time)
}

type AssessmentRepo interface {
	// ใช้จาก repo เดิมของคุณ
	GetAssessmentByID(id string) (*models.Assessment, error)
}
