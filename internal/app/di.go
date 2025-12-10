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

	// assessment
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
	ContentHTTP      *contenthandler.Handler
	AssessHTTP       *assesshandler.Handler
	AttemptHTTP      *assesshandler.AttemptHandler
	LearningHTTP     *learnhdl.Handler
	CourseHTTP       *contenthandler.CourseHandler
	CategoryHTTP     *contenthandler.CategoryHandler
	MyHandler        *learnhdl.MyHandler
	AnalyticsHandler *learnhdl.AnalyticsHandler // <<< เพิ่มตรงนี้
}

func Build() Deps {
	// ===== Users/Auth =====
	ur := usersrepo.NewGormRepository(config.DB)
	us := usersvc.NewService(ur)
	ar := authrepo.NewGormRepository(config.DB)
	as := authsvc.New(ur, ar)

	// ===== Content =====
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

	// ===== Assessment =====
	assRepo := assessrepo.New(config.DB)
	asSvc := assesssvc.New(assRepo)
	assHTTP := assesshandler.New(asSvc)

	// attempt repo (assessment attempts)
	attRepo := assessrepo.NewAttemptRepo(config.DB)

	// ===== Learning =====
	lr := learnrepo.New(config.DB)

	// MyCourses (for /my endpoints)
	myCoursesRepo := learnrepo.NewMyCoursesRepo(config.DB)
	myCoursesSvc := learnsvc.NewMyCoursesService(myCoursesRepo)
	myHandler := learnhdl.NewMyHandler(myCoursesSvc)

	// Metrics (analytics)
	metricsRepo := learnrepo.NewMetricsRepo(config.DB)
	metricsSvc := learnsvc.NewMetricsService(metricsRepo)
	analyticsHandler := learnhdl.NewAnalyticsHandler(metricsSvc)

	// Learning service (main)
	ls := learnsvc.New(lr, metricsSvc)
	lh := learnhdl.New(ls)

	// Attempt service (assessment attempts) — ปรับตาม signature ของคุณ
	// ถ้า NewAttemptService ต้องการ (assRepo, attRepo) เป็นอันพอ
	// attSvc := assesssvc.NewAttemptService(assRepo, attRepo)
	attSvc := assesssvc.NewAttemptService(assRepo, attRepo, lr, metricsSvc)
	attHTTP := assesshandler.NewAttemptHandler(attSvc)

	// Courses/Category
	courseRepo := contentrepo.NewCourseRepo(config.DB)
	moduleRepo := contentrepo.NewModuleRepo(config.DB)
	categoryRepo := contentrepo.NewCategoryRepo(config.DB)
	courseDeptRepo := contentrepo.NewCourseDeptRepo(config.DB)
	courseSvc := contentservice.NewCourseService(courseRepo, moduleRepo, lessonRepo, categoryRepo, courseDeptRepo)
	categorySvc := contentservice.NewCategoryService(categoryRepo, courseRepo)
	courseHTTP := contenthandler.NewCourseHandler(courseSvc)
	categoryHTTP := contenthandler.NewCategoryHandler(categorySvc)

	return Deps{
		UserSvc:          us,
		AuthSvc:          as,
		ContentHTTP:      contentHTTP,
		AssessHTTP:       assHTTP,
		AttemptHTTP:      attHTTP,
		LearningHTTP:     lh,
		CourseHTTP:       courseHTTP,
		CategoryHTTP:     categoryHTTP,
		MyHandler:        myHandler,
		AnalyticsHandler: analyticsHandler,
	}
}
