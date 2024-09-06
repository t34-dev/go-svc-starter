// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: auth_v1/auth.proto

package auth_v1

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

// Validate checks the field values on LoginRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoginRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoginRequestMultiError, or
// nil if none found.
func (m *LoginRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if l := utf8.RuneCountInString(m.GetUsername()); l < 1 || l > 50 {
		err := LoginRequestValidationError{
			field:  "Username",
			reason: "value length must be between 1 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if l := utf8.RuneCountInString(m.GetPassword()); l < 1 || l > 50 {
		err := LoginRequestValidationError{
			field:  "Password",
			reason: "value length must be between 1 and 50 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return LoginRequestMultiError(errors)
	}

	return nil
}

// LoginRequestMultiError is an error wrapping multiple validation errors
// returned by LoginRequest.ValidateAll() if the designated constraints aren't met.
type LoginRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginRequestMultiError) AllErrors() []error { return m }

// LoginRequestValidationError is the validation error returned by
// LoginRequest.Validate if the designated constraints aren't met.
type LoginRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginRequestValidationError) ErrorName() string { return "LoginRequestValidationError" }

// Error satisfies the builtin error interface
func (e LoginRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginRequestValidationError{}

// Validate checks the field values on LoginResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoginResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoginResponseMultiError, or
// nil if none found.
func (m *LoginResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccessToken

	// no validation rules for RefreshToken

	if len(errors) > 0 {
		return LoginResponseMultiError(errors)
	}

	return nil
}

// LoginResponseMultiError is an error wrapping multiple validation errors
// returned by LoginResponse.ValidateAll() if the designated constraints
// aren't met.
type LoginResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginResponseMultiError) AllErrors() []error { return m }

// LoginResponseValidationError is the validation error returned by
// LoginResponse.Validate if the designated constraints aren't met.
type LoginResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginResponseValidationError) ErrorName() string { return "LoginResponseValidationError" }

// Error satisfies the builtin error interface
func (e LoginResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginResponseValidationError{}

// Validate checks the field values on UpdateTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateTokenRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateTokenRequestMultiError, or nil if none found.
func (m *UpdateTokenRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateTokenRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for RefreshToken

	if len(errors) > 0 {
		return UpdateTokenRequestMultiError(errors)
	}

	return nil
}

// UpdateTokenRequestMultiError is an error wrapping multiple validation errors
// returned by UpdateTokenRequest.ValidateAll() if the designated constraints
// aren't met.
type UpdateTokenRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateTokenRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateTokenRequestMultiError) AllErrors() []error { return m }

// UpdateTokenRequestValidationError is the validation error returned by
// UpdateTokenRequest.Validate if the designated constraints aren't met.
type UpdateTokenRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateTokenRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateTokenRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateTokenRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateTokenRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateTokenRequestValidationError) ErrorName() string {
	return "UpdateTokenRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateTokenRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateTokenRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateTokenRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateTokenRequestValidationError{}

// Validate checks the field values on UpdateTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateTokenResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateTokenResponseMultiError, or nil if none found.
func (m *UpdateTokenResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateTokenResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccessToken

	// no validation rules for RefreshToken

	if len(errors) > 0 {
		return UpdateTokenResponseMultiError(errors)
	}

	return nil
}

// UpdateTokenResponseMultiError is an error wrapping multiple validation
// errors returned by UpdateTokenResponse.ValidateAll() if the designated
// constraints aren't met.
type UpdateTokenResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateTokenResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateTokenResponseMultiError) AllErrors() []error { return m }

// UpdateTokenResponseValidationError is the validation error returned by
// UpdateTokenResponse.Validate if the designated constraints aren't met.
type UpdateTokenResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateTokenResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateTokenResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateTokenResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateTokenResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateTokenResponseValidationError) ErrorName() string {
	return "UpdateTokenResponseValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateTokenResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateTokenResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateTokenResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateTokenResponseValidationError{}
