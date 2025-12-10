package app

import (
	"os"

	authhandler "github.com/Marugo/birdlax/internal/modules/auth/handler"
	userhandler "github.com/Marugo/birdlax/internal/modules/user/handler"

	assesshandler "github.com/Marugo/birdlax/internal/modules/assessment/handler"
	contenthandler "github.com/Marugo/birdlax/internal/modules/content/handler"
	learninghandler "github.com/Marugo/birdlax/internal/modules/learning/handler"

	"github.com/Marugo/birdlax/internal/shared/middleware"
	"github.com/gofiber/fiber/v2"
	fibercors "github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Register(app *fiber.App) {
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(middleware.RequestID())

	// 1) CORS: อนุญาต frontend ที่ :8000 + เฮดเดอร์สำหรับวิดีโอ
	origins := os.Getenv("APP_CORS_ORIGINS")
	if origins == "" {
		origins = "http://localhost:8000"
	}
	app.Use(fibercors.New(fibercors.Config{
		AllowOrigins:     origins, // e.g. http://localhost:8000
		AllowMethods:     "GET,HEAD,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Range",
		ExposeHeaders:    "Content-Length, Content-Range, Accept-Ranges",
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// 2) Static videos: รองรับ Byte Range และใส่ Accept-Ranges ให้ชัด
	publicBaseURL := os.Getenv("PUBLIC_BASE_URL")
	if publicBaseURL == "" {
		publicBaseURL = "/static/videos"
	}
	uploadBaseDir := os.Getenv("UPLOAD_BASE_DIR")
	if uploadBaseDir == "" {
		uploadBaseDir = "/data/uploads/videos"
	}

	// เสิร์ฟไฟล์แบบ byte range (จำเป็นต่อการ seek/สตรีม)
	app.Static(publicBaseURL, uploadBaseDir, fiber.Static{
		Compress:  true,
		ByteRange: true, // ✅ สำคัญสำหรับ <video> seek
		Browse:    false,
	})

	// เติม header Accept-Ranges ให้ทุกคำขอใต้ /static/videos
	app.Use(publicBaseURL+"/*", func(c *fiber.Ctx) error {
		c.Set("Accept-Ranges", "bytes")                                 // ✅ video
		c.Set("Cache-Control", "public, max-age=3600, must-revalidate") // ✅ แทน CacheControl ใน struct
		// c.Set("Cross-Origin-Resource-Policy", "cross-origin")            // (ถ้าต้องการ)
		return c.Next()
	})

	// (ไม่จำเป็นเสมอไป เพราะ CORS middleware จัดการแล้ว)
	// app.Options("/*", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	// 3) health
	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendString("ok") })

	// 4) DI
	deps := Build()

	// 5) API v1
	api := app.Group("/api/v1")

	// Public
	authHTTP := authhandler.NewHTTPHandler(deps.AuthSvc, deps.UserSvc)
	authhandler.Register(api, authHTTP)

	// Protected
	protected := api.Group("", middleware.AuthRequired())

	userhandler.Register(protected, deps.UserSvc)
	contenthandler.Register(protected, deps.ContentHTTP)
	assesshandler.Register(protected, deps.AssessHTTP)
	learninghandler.Register(protected, deps.LearningHTTP)
	contenthandler.RegisterCourseRoutes(protected, deps.CourseHTTP)
	assesshandler.RegisterAttemptRoutes(protected, deps.AttemptHTTP)
	contenthandler.RegisterCategoryRoutes(api, deps.CategoryHTTP)
	learninghandler.MyRegister(protected, deps.LearningHTTP, deps.MyHandler)
	learninghandler.RegisterAdminRoutes(protected, deps.AnalyticsHandler)

}
