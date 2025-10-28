package response

import "github.com/gofiber/fiber/v2"

type Meta struct {
	RequestID string `json:"request_id,omitempty"`
}

type Resp struct {
	Code    string      `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

func JSON(ctx *fiber.Ctx, status int, code, msg string, data interface{}) error {
	reqID := ctx.GetRespHeader("X-Request-ID")
	return ctx.Status(status).JSON(Resp{
		Code:    code,
		Message: msg,
		Data:    data,
		Meta:    &Meta{RequestID: reqID},
	})
}

func OK(ctx *fiber.Ctx, data interface{}) error {
	return JSON(ctx, fiber.StatusOK, "OK", "", data)
}

func Err(ctx *fiber.Ctx, status int, code, msg string) error {
	return JSON(ctx, status, code, msg, nil)
}
