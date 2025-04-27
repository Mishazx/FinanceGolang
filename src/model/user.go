package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Fio      string `json:"fio" gorm:"not null"`
	Username string `json:"username" gorm:"unique;not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"`
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
