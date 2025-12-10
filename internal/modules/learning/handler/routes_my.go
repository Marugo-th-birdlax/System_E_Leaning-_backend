package handler

import "github.com/gofiber/fiber/v2"

func MyRegister(r fiber.Router, _ *Handler, my *MyHandler) {
	myGroup := r.Group("/my") // ใช้ protected group มาจากข้างนอกแล้ว ไม่ต้อง Auth อีก
	myGroup.Get("/department-courses", my.MyDepartmentCourses)
	myGroup.Get("/courses", my.MyCourses)

	myGroup.Get("/courses/:courseID/progress", my.MyCourseProgress)

}
