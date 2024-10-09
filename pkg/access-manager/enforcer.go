package access_manager

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	stringadapter "github.com/casbin/casbin/v2/persist/string-adapter"
	"strings"
)

func createEnforcer(modelStr, policyStr string) (*casbin.Enforcer, error) {
	// Check model validity
	m, err := model.NewModelFromString(modelStr)
	if err != nil {
		return nil, fmt.Errorf("invalid model data: %w", err)
	}

	// Create adapter from policy data
	err = validatePolicyUsingModelAddDef(modelStr, policyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid policy data: %w", err)
	}
	sa := stringadapter.NewAdapter(policyStr)

	// Create enforcer
	e, err := casbin.NewEnforcer(m, sa)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	// Attempt to load policy
	err = e.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("invalid policy data: %w", err)
	}

	return e, nil
}
func validatePolicyUsingModelAddDef(modelStr, policyStr string) error {
	m, err := model.NewModelFromString(modelStr)
	if err != nil {
		return fmt.Errorf("invalid model: %w", err)
	}

	for _, line := range strings.Split(policyStr, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.Split(line, ",")
		for i := range tokens {
			tokens[i] = strings.TrimSpace(tokens[i])
		}

		if len(tokens) < 2 {
			return fmt.Errorf("invalid policy line: %s", line)
		}

		ptype := tokens[0]
		if !m.AddDef(ptype, ptype, strings.Join(tokens[1:], ", ")) {
			return fmt.Errorf("invalid policy: %w", err)
		}
	}

	return nil
}
