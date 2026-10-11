package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

type myStruct struct {
	in          interface{}
	expectedErr error
}

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []myStruct{
		{
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John Doe",
				Age:    30,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John Doe",
				Age:    16,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: fmt.Errorf("value 16 is less than minimum 18")},
			},
		},
		{
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John Doe",
				Age:    30,
				Email:  "john.doe",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Email", Err: fmt.Errorf("value john.doe does not match regexp ^\\w+@\\w+\\.\\w+$")},
			},
		},
		{
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John Doe",
				Age:    30,
				Email:  "john@example.com",
				Role:   "user",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Role", Err: fmt.Errorf("value user is not in allowed set [admin stuff]")},
			},
		},
		{
			in: User{
				ID:     "123e4567-e89b-12d3-a456-426614174000",
				Name:   "John Doe",
				Age:    30,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"1234567890"},
			},
			expectedErr: ValidationErrors{
				{Field: "Phones", Err: fmt.Errorf("value 1234567890 does not have length 11")},
			},
		},
		{
			in: App{
				Version: "v1.01",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "v1",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("value v1 does not have length 5")},
			},
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 403,
				Body: "Forbidden",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("value 403 is not in allowed set [200 404 500]")},
			},
		},
		{
			in:          42,
			expectedErr: fmt.Errorf("input must be a struct"),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error %v, got no error", tt.expectedErr)
				}
				validate(t, tt, err)
			}
		})
	}
}

func validate(t *testing.T, tt myStruct, err error) {
	t.Helper()
	var validationErrors ValidationErrors
	if errors.As(tt.expectedErr, &validationErrors) {
		var actualErrors ValidationErrors
		ok := errors.As(err, &actualErrors)
		if !ok {
			t.Errorf("Expected ValidationErrors, got %T", err)
		}
		if len(actualErrors) != len(validationErrors) {
			t.Errorf("Expected %d validation errors, got %d", len(validationErrors), len(actualErrors))
		}
		for i, expectedErr := range validationErrors {
			if actualErrors[i].Field != expectedErr.Field || actualErrors[i].Err.Error() != expectedErr.Err.Error() {
				t.Errorf("Expected error %v, got %v", expectedErr, actualErrors[i])
			}
		}
	}
}
