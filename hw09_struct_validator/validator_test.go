package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

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
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:    "12345678-1234-5678-1234-567812345678",
				Name:  "Alice",
				Age:   30,
				Email: "alice@example.com",
				Role:  "admin",
				Phones: []string{
					"12345678901",
					"23456789012",
				},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:    "short-id",
				Name:  "Bob",
				Age:   17,
				Email: "invalid-email",
				Role:  "user",
				Phones: []string{
					"1234567",
				},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("length must be 36")},
				{Field: "Age", Err: fmt.Errorf("value must be at least 18")},
				{Field: "Email", Err: fmt.Errorf("value does not match regexp: ^\\w+@\\w+\\.\\w+$")},
				{Field: "Role", Err: fmt.Errorf("value must be one of [admin stuff]")},
				{Field: "Phones", Err: fmt.Errorf("length must be 11")},
			},
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.0",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("length must be 5")},
			},
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
				Code: 201,
				Body: "Created",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("value must be one of [200 404 500]")},
			},
		},
		{
			in: struct {
				Field string `validate:"len"`
			}{
				Field: "value",
			},
			expectedErr: fmt.Errorf("field Field: %w", ErrInvalidValidatorFormat),
		},
		{
			in: struct {
				Field int `validate:"unknown:10"`
			}{
				Field: 5,
			},
			expectedErr: fmt.Errorf("field Field: %w", ErrInvalidValidatorType),
		},
		{
			in: struct {
				Field string `validate:"regexp:[invalid"`
			}{
				Field: "value",
			},
			expectedErr: fmt.Errorf("field Field: %w", ErrInvalidRegexp),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
