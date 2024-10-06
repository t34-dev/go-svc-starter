package model

import (
	"strings"
)

type RoleX string

const (
	AdminRole RoleX = "Admin"
	VipRole   RoleX = "Vip"
	// additional roles as needed
)

type RolesX []RoleX

// String Convert roles to a comma-separated string
func (r RolesX) String() string {
	rolesAsStrings := make([]string, len(r))
	for i, role := range r {
		rolesAsStrings[i] = string(role)
	}
	return strings.Join(rolesAsStrings, ",")
}

// HasRoles Checks if all specified roles are present in the current array of roles
func (r RolesX) HasRoles(roles ...RoleX) bool {
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
func (r *RolesX) AddRoles(roles ...RoleX) {
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
func (r *RolesX) RemoveRoles(roles ...RoleX) {
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
