package handler

import (
	"github.com/Marugo/birdlax/internal/modules/user"
	"github.com/gofiber/fiber/v2"
)

func Register(r fiber.Router, svc user.Service) {
	h := NewHTTPHandler(svc)
	g := r.Group("/users")
	g.Get("/", h.List)
	g.Post("/", h.Create)
	g.Get("/:id", h.Get)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)
}
