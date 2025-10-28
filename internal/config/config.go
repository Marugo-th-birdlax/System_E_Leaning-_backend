package config

import (
	"fmt"
	"log"
	"os"
	"time"

	assessmentmodels "github.com/Marugo/birdlax/internal/modules/assessment/models"
	contentmodels "github.com/Marugo/birdlax/internal/modules/content/models"
	learningmodels "github.com/Marugo/birdlax/internal/modules/learning/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	authrepo "github.com/Marugo/birdlax/internal/modules/auth/repo"
	usermodels "github.com/Marugo/birdlax/internal/modules/user/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
)

var DB *gorm.DB

func Init() error {
	_ = godotenv.Load()

	locName := getEnv("APP_TIMEZONE", "Asia/Bangkok")
	loc, err := time.LoadLocation(locName)
	if err != nil {
		return fmt.Errorf("load tz: %w", err)
	}
	time.Local = loc

	if err := ConnectDatabase(); err != nil {
		return err
	}

	// Auto-migrate เฉพาะโมดูล user
	if err := DB.AutoMigrate(
		&usermodels.User{},
		&authrepo.RefreshToken{},
		&contentmodels.Asset{},
		&contentmodels.Lesson{},
		&assessmentmodels.Assessment{},
		&assessmentmodels.Question{},
		&assessmentmodels.Choice{},
		&learningmodels.Enrollment{},
		&learningmodels.UserLessonProgress{},
		&contentmodels.Course{},
		&contentmodels.CourseModule{},
		&assessmentmodels.Attempt{},
		&assessmentmodels.Answer{},
	); err != nil {
		return fmt.Errorf("auto-migrate: %w", err)
	}

	return nil
}

func AppName() string { return getEnv("APP_NAME", "go-fiber-gorm") }
func AppPort() string { return getEnv("APP_PORT", "3000") }

func ConnectDatabase() error {
	driver := getEnv("DB_DRIVER", "mysql")

	switch driver {
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			getEnv("DB_USER", "root"),
			getEnv("DB_PASS", ""),
			getEnv("DB_HOST", "127.0.0.1"),
			getEnv("DB_PORT", "3306"),
			getEnv("DB_NAME", "appdb"),
		)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("db connect: %w", err)
		}
		DB = db

	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_USER", "postgres"),
			getEnv("DB_PASS", "postgres"),
			getEnv("DB_NAME", "appdb"),
			getEnv("DB_PORT", "5432"),
			getEnv("DB_SSLMODE", "disable"),
			getEnv("APP_TIMEZONE", "Asia/Bangkok"),
		)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("db connect: %w", err)
		}
		DB = db

	default:
		return fmt.Errorf("unsupported DB_DRIVER: %s", driver)
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(2 * time.Hour)

	log.Println("database connected")
	return nil
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
