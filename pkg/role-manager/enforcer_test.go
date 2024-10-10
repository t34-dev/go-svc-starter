package role_manager

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEnforcer(t *testing.T) {
	tests := []struct {
		name      string
		modelStr  string
		policyStr string
		wantErr   bool
		errMsg    string
	}{
		{
			name: "Valid model and policy",
			modelStr: `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`,
			policyStr: `
p, alice, data1, read
p, bob, data2, write
`,
			wantErr: false,
		},
		{
			name:      "Invalid model",
			modelStr:  "invalid model",
			policyStr: "p, alice, data1, read",
			wantErr:   true,
			errMsg:    "invalid model data",
		},
		{
			name: "Valid model with invalid policy",
			modelStr: `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`,
			policyStr: "invalid policy",
			wantErr:   true,
			errMsg:    "invalid policy data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := createEnforcer(tt.modelStr, tt.policyStr)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, e)
				assert.True(t, strings.Contains(err.Error(), tt.errMsg))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, e)
				// Additional checks for the created enforcer
				ok, err := e.Enforce("alice", "data1", "read")
				assert.NoError(t, err)
				assert.True(t, ok)
				ok, err = e.Enforce("bob", "data2", "write")
				assert.NoError(t, err)
				assert.True(t, ok)
				ok, err = e.Enforce("alice", "data2", "write")
				assert.NoError(t, err)
				assert.False(t, ok)
			}
		})
	}
}
