package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/Marugo/birdlax/internal/app"
	"github.com/Marugo/birdlax/internal/config"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("config init error: %v", err)
	}

	appSrv := fiber.New(fiber.Config{
		AppName:      config.AppName(),
		ServerHeader: "fiber",
		BodyLimit:    200 * 1024 * 1024,
	})

	app.Register(appSrv)

	addr := ":" + config.AppPort()
	// If you want to print registered routes, use Fiber's built-in method:
	for _, route := range appSrv.GetRoutes() {
		println(route.Method, route.Path)
	}

	log.Printf("listening on %s", addr)
	if err := appSrv.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
