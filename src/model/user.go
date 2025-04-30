package model

import (
	"regexp"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Fio      string `json:"fio" gorm:"not null"`
	Username string `json:"username" gorm:"unique;not null"`
	Email    string `json:"email" gorm:"unique;not null"` // Removed CHECK constraint
	Password string `json:"-" gorm:"not null"`            // Removed CHECK constraint
	Roles    []Role `json:"roles" gorm:"many2many:user_roles;"`
}

func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

func (u *User) IsAdmin() bool {
	return u.HasRole(RoleAdmin)
}

// ValidateEmail проверяет корректность email
func (u *User) ValidateEmail() bool {
	emailRegex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	return emailRegex.MatchString(u.Email)
}

// ValidatePassword проверяет корректность пароля
func (u *User) ValidatePassword() bool {
	return len(u.Password) >= 8
}
