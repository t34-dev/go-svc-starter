package model

import (
	"strings"
)

type Role string

const (
	AdminRole Role = "Admin"
	VipRole   Role = "Vip"
	// дополнительные роли по мере необходимости
)

type Roles []Role

// String Преобразовать роли в строку через запятую
func (r Roles) String() string {
	rolesAsStrings := make([]string, len(r))
	for i, role := range r {
		rolesAsStrings[i] = string(role)
	}
	return strings.Join(rolesAsStrings, ",")
}

// HasRoles Проверяет наличие всех указанных ролей в текущем массиве ролей
func (r Roles) HasRoles(roles ...Role) bool {
	for _, roleToCheck := range roles {
		found := false
		for _, existingRole := range r {
			if roleToCheck == existingRole {
				found = true
				break
			}
		}
		// Если хотя бы одна из ролей не найдена, вернем false
		if !found {
			return false
		}
	}
	return true
}

// AddRoles добавляет роли, которые еще не присутствуют
func (r *Roles) AddRoles(roles ...Role) {
	for _, roleToAdd := range roles {
		found := false
		for _, existingRole := range *r {
			if roleToAdd == existingRole {
				found = true
				break
			}
		}
		if !found {
			*r = append(*r, roleToAdd)
		}
	}
}

// RemoveRoles удаляет роли, если они присутствуют
func (r *Roles) RemoveRoles(roles ...Role) {
	for _, roleToRemove := range roles {
		indexToRemove := -1
		for i, existingRole := range *r {
			if roleToRemove == existingRole {
				indexToRemove = i
				break
			}
		}
		if indexToRemove != -1 {
			*r = append((*r)[:indexToRemove], (*r)[indexToRemove+1:]...)
		}
	}
}
