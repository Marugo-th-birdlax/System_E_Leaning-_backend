package app

import (
	"os"

	"github.com/Marugo/birdlax/internal/config"

	// auth & users
	"github.com/Marugo/birdlax/internal/modules/auth"
	authrepo "github.com/Marugo/birdlax/internal/modules/auth/repo"
	authsvc "github.com/Marugo/birdlax/internal/modules/auth/service"
	"github.com/Marugo/birdlax/internal/modules/user"
	usersrepo "github.com/Marugo/birdlax/internal/modules/user/repo"
	usersvc "github.com/Marugo/birdlax/internal/modules/user/service"

	// content
	contenthandler "github.com/Marugo/birdlax/internal/modules/content/handler"
	contentrepo "github.com/Marugo/birdlax/internal/modules/content/repo"
	contentservice "github.com/Marugo/birdlax/internal/modules/content/service"
	contentstorage "github.com/Marugo/birdlax/internal/modules/content/storage"

	// assessment (ใช้ alias ชุดเดียว)
	assesshandler "github.com/Marugo/birdlax/internal/modules/assessment/handler"
	assessrepo "github.com/Marugo/birdlax/internal/modules/assessment/repo"
	assesssvc "github.com/Marugo/birdlax/internal/modules/assessment/service"

	// learning
	learnhdl "github.com/Marugo/birdlax/internal/modules/learning/handler"
	learnrepo "github.com/Marugo/birdlax/internal/modules/learning/repo"
	learnsvc "github.com/Marugo/birdlax/internal/modules/learning/service"
)

type Deps struct {
	UserSvc user.Service
	AuthSvc auth.Service

	// HTTP handlers
	ContentHTTP  *contenthandler.Handler
	AssessHTTP   *assesshandler.Handler
	AttemptHTTP  *assesshandler.AttemptHandler // ✅ เพิ่มฟิลด์นี้
	LearningHTTP *learnhdl.Handler
	CourseHTTP   *contenthandler.CourseHandler
}

func Build() Deps {
	// ===== Users/Auth =====
	ur := usersrepo.NewGormRepository(config.DB)
	us := usersvc.NewService(ur)

	ar := authrepo.NewGormRepository(config.DB)
	as := authsvc.New(ur, ar)

	// ===== Content (upload video / lesson) =====
	uploadBaseDir := os.Getenv("UPLOAD_BASE_DIR")
	if uploadBaseDir == "" {
		uploadBaseDir = "/data/uploads/videos"
	}
	publicBaseURL := os.Getenv("PUBLIC_BASE_URL")
	if publicBaseURL == "" {
		publicBaseURL = "/static/videos"
	}
	uploader := &contentstorage.LocalFS{BaseDir: uploadBaseDir, BaseURL: publicBaseURL}

	assetRepo := contentrepo.NewAssetRepo(config.DB)
	lessonRepo := contentrepo.NewLessonRepo(config.DB)
	contentSvc := contentservice.New(assetRepo, lessonRepo, uploader)
	contentHTTP := contenthandler.New(contentSvc)

	// ===== Assessment (create / add question) =====
	assRepo := assessrepo.New(config.DB)
	asSvc := assesssvc.New(assRepo)
	assHTTP := assesshandler.New(asSvc)

	// ===== Attempts (start/answer/submit) =====
	attRepo := assessrepo.NewAttemptRepo(config.DB)
	attSvc := assesssvc.NewAttemptService(assRepo, attRepo)
	attHTTP := assesshandler.NewAttemptHandler(attSvc)

	// ===== Learning =====
	lr := learnrepo.New(config.DB)
	ls := learnsvc.New(lr)
	lh := learnhdl.New(ls)

	// ===== Courses/Modules =====
	courseRepo := contentrepo.NewCourseRepo(config.DB)
	moduleRepo := contentrepo.NewModuleRepo(config.DB)
	courseSvc := contentservice.NewCourseService(courseRepo, moduleRepo, lessonRepo)
	courseHTTP := contenthandler.NewCourseHandler(courseSvc)

	return Deps{
		UserSvc:      us,
		AuthSvc:      as,
		ContentHTTP:  contentHTTP,
		AssessHTTP:   assHTTP,
		AttemptHTTP:  attHTTP, // ✅
		LearningHTTP: lh,
		CourseHTTP:   courseHTTP,
	}
}
