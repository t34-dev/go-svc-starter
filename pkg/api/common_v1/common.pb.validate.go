// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: common_v1/common.proto

package common_v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on TimeResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *TimeResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TimeResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in TimeResponseMultiError, or
// nil if none found.
func (m *TimeResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *TimeResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetTime()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TimeResponseValidationError{
					field:  "Time",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TimeResponseValidationError{
					field:  "Time",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTime()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TimeResponseValidationError{
				field:  "Time",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return TimeResponseMultiError(errors)
	}

	return nil
}

// TimeResponseMultiError is an error wrapping multiple validation errors
// returned by TimeResponse.ValidateAll() if the designated constraints aren't met.
type TimeResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TimeResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TimeResponseMultiError) AllErrors() []error { return m }

// TimeResponseValidationError is the validation error returned by
// TimeResponse.Validate if the designated constraints aren't met.
type TimeResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TimeResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TimeResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TimeResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TimeResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TimeResponseValidationError) ErrorName() string { return "TimeResponseValidationError" }

// Error satisfies the builtin error interface
func (e TimeResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTimeResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TimeResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TimeResponseValidationError{}

// Validate checks the field values on LongOperationRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *LongOperationRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LongOperationRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// LongOperationRequestMultiError, or nil if none found.
func (m *LongOperationRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *LongOperationRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return LongOperationRequestMultiError(errors)
	}

	return nil
}

// LongOperationRequestMultiError is an error wrapping multiple validation
// errors returned by LongOperationRequest.ValidateAll() if the designated
// constraints aren't met.
type LongOperationRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LongOperationRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LongOperationRequestMultiError) AllErrors() []error { return m }

// LongOperationRequestValidationError is the validation error returned by
// LongOperationRequest.Validate if the designated constraints aren't met.
type LongOperationRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LongOperationRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LongOperationRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LongOperationRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LongOperationRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LongOperationRequestValidationError) ErrorName() string {
	return "LongOperationRequestValidationError"
}

// Error satisfies the builtin error interface
func (e LongOperationRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLongOperationRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LongOperationRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LongOperationRequestValidationError{}

// Validate checks the field values on LongOperationResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *LongOperationResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LongOperationResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// LongOperationResponseMultiError, or nil if none found.
func (m *LongOperationResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *LongOperationResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Status

	// no validation rules for Message

	// no validation rules for Progress

	// no validation rules for Result

	if len(errors) > 0 {
		return LongOperationResponseMultiError(errors)
	}

	return nil
}

// LongOperationResponseMultiError is an error wrapping multiple validation
// errors returned by LongOperationResponse.ValidateAll() if the designated
// constraints aren't met.
type LongOperationResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LongOperationResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LongOperationResponseMultiError) AllErrors() []error { return m }

// LongOperationResponseValidationError is the validation error returned by
// LongOperationResponse.Validate if the designated constraints aren't met.
type LongOperationResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LongOperationResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LongOperationResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LongOperationResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LongOperationResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LongOperationResponseValidationError) ErrorName() string {
	return "LongOperationResponseValidationError"
}

// Error satisfies the builtin error interface
func (e LongOperationResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLongOperationResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LongOperationResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LongOperationResponseValidationError{}
