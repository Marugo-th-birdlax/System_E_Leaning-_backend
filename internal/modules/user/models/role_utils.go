// internal/modules/user/models/role_utils.go
package models

import "strings"

// Validate ว่าค่านี้อยู่ในเซ็ต roles เราหรือเปล่า
func IsValidRole(r Role) bool {
	for _, x := range AllRoles() {
		if x == r {
			return true
		}
	}
	return false
}

// ช่วย parse จาก string ที่มากับ request (case-insensitive)
func ParseRole(s string) (Role, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	r := Role(s)
	return r, IsValidRole(r)
}

// กำหนด “แรงดีกรีสิทธิ์” ไว้เทียบขั้นต่ำ (อยากปรับสเกลได้)
var roleRank = map[Role]int{
	RoleEmployee: 1,
	RoleHR:       2,
	RoleAdmin:    3,
}

// ใช้สำหรับเช็กว่า role ผู้ใช้อย่างน้อยต้อง = min ขึ้นไป
func IsAtLeast(r Role, min Role) bool {
	return roleRank[r] >= roleRank[min]
}

func IsValidDepartmentRole(r DepartmentRole) bool {
	for _, x := range AllDepartmentRoles() {
		if x == r {
			return true
		}
	}
	return false
}

func ParseDepartmentRole(s string) (DepartmentRole, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	r := DepartmentRole(s)
	return r, IsValidDepartmentRole(r)
}
