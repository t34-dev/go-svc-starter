package model

import (
	"strings"
)

type Role string

const (
	AdminRole Role = "Admin"
	VipRole   Role = "Vip"
	// additional roles as needed
)

type Roles []Role

// String Convert roles to a comma-separated string
func (r Roles) String() string {
	rolesAsStrings := make([]string, len(r))
	for i, role := range r {
		rolesAsStrings[i] = string(role)
	}
	return strings.Join(rolesAsStrings, ",")
}

// HasRoles Checks if all specified roles are present in the current array of roles
func (r Roles) HasRoles(roles ...Role) bool {
	for _, roleToCheck := range roles {
		found := false
		for _, existingRole := range r {
			if roleToCheck == existingRole {
				found = true
				break
			}
		}
		// If at least one of the roles is not found, return false
		if !found {
			return false
		}
	}
	return true
}

// AddRoles adds roles that are not already present
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

// RemoveRoles removes roles if they are present
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
