package validator

import (
	"context"
	"github.com/t34-dev/go-utils/pkg/sys/validate"
)

func ValidateID(id int64) validate.Condition {
	return func(ctx context.Context) error {
		errors := []string{}
		if id <= 0 {
			errors = append(errors, "id must be greater than 0")
		}
		if id > 100 {
			errors = append(errors, "id must be less than 100")
		}

		if len(errors) > 0 {
			return validate.NewValidationErrors(errors...)
		}

		return nil
	}
}
