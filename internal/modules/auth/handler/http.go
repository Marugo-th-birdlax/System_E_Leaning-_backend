package handler

import (
	"time"

	auth "github.com/Marugo/birdlax/internal/modules/auth"
	user "github.com/Marugo/birdlax/internal/modules/user"
	"github.com/Marugo/birdlax/internal/modules/user/dto"
	"github.com/Marugo/birdlax/internal/shared/response"
	"github.com/gofiber/fiber/v2"
)

type HTTPHandler struct {
	svc     auth.Service
	userSvc user.Service
}

func NewHTTPHandler(s auth.Service, u user.Service) *HTTPHandler {
	return &HTTPHandler{svc: s, userSvc: u}
}

type loginReq struct {
	EmployeeCode string `json:"employee_code"`
	Password     string `json:"password"`
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *HTTPHandler) Login(c *fiber.Ctx) error {
	var r loginReq
	if err := c.BodyParser(&r); err != nil {
		return response.Err(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid JSON")
	}
	access, refresh, err := h.svc.Login(c.Context(), r.EmployeeCode, r.Password) // 👈 ส่ง employee_code
	if err != nil {
		return response.Err(c, fiber.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
	}
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh, // refresh token string
		HTTPOnly: true,
		Secure:   false, // true ถ้าใช้ https
		SameSite: "Lax", // fiber v2 ใช้ fiber.CookieSameSiteLaxMode ได้เช่นกัน
		Path:     "/",
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	})
	return response.OK(c, fiber.Map{
		"access_token": access, // ให้ frontend ใส่ลง Authorization header
	})
}

func (h *HTTPHandler) Refresh(c *fiber.Ctx) error {
	// รับจาก body หรือ cookie ก็ได้ — แนะนำจาก cookie
	rt := c.Cookies("refresh_token")
	if rt == "" {
		var r struct {
			RefreshToken string `json:"refresh_token"`
		}
		_ = c.BodyParser(&r)
		rt = r.RefreshToken
	}
	access, newRefresh, err := h.svc.Refresh(c.Context(), rt)
	if err != nil {
		return response.Err(c, fiber.StatusUnauthorized, "INVALID_REFRESH", err.Error())
	}

	// rotate cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefresh,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	})
	return response.OK(c, fiber.Map{"access_token": access})
}

func (h *HTTPHandler) Logout(c *fiber.Ctx) error {
	rt := c.Cookies("refresh_token")
	if rt == "" {
		var r struct {
			RefreshToken string `json:"refresh_token"`
		}
		_ = c.BodyParser(&r)
		rt = r.RefreshToken
	}
	_ = h.svc.Logout(c.Context(), rt)
	// clear cookie
	c.Cookie(&fiber.Cookie{
		Name:   "refresh_token",
		Value:  "",
		MaxAge: -1, Path: "/",
		HTTPOnly: true, Secure: false, SameSite: "Lax",
	})
	return response.OK(c, fiber.Map{"ok": true})
}

func (h *HTTPHandler) Me(c *fiber.Ctx) error {
	userID, _ := c.Locals("user_id").(string)
	role, _ := c.Locals("role").(string)

	u, err := h.userSvc.Get(c.Context(), userID)
	if err != nil || u == nil {
		return response.Err(c, fiber.StatusNotFound, "NOT_FOUND", "user not found")
	}

	out := dto.FromModel(u) // มี first_name/last_name/ฯลฯ อยู่แล้ว
	// role ใน DB อาจต่างจาก claim (กรณีเพิ่งเปลี่ยนสิทธิ์) — เลือกหนึ่ง:
	// 1) เชื่อ DB (ค่าจริงปัจจุบัน):
	//    out.Role = string(u.Role)
	// 2) หรือเชื่อ Token:
	_ = role // ใช้ถ้าต้องการยึดตาม token

	// ตอบเฉพาะฟิลด์ที่อยากโชว์ (ถ้าต้องการลดรูป)
	return response.OK(c, fiber.Map{
		"id":            out.ID,
		"employee_code": out.EmployeeCode,
		"email":         out.Email,
		"first_name":    out.FirstName,
		"last_name":     out.LastName,
		"role":          out.Role, // หรือ role จาก token
		"phone":         out.Phone,
		"is_active":     out.IsActive,
	})
}
