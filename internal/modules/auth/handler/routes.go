package handler

import (
	"github.com/Marugo/birdlax/internal/shared/middleware"
	"github.com/gofiber/fiber/v2"
)

func Register(r fiber.Router, h *HTTPHandler) {
	g := r.Group("/auth")
	g.Post("/login", h.Login)
	g.Post("/refresh", h.Refresh)
	g.Post("/logout", middleware.AuthRequired(), h.Logout)
	g.Get("/me", middleware.AuthRequired(), h.Me)
}
