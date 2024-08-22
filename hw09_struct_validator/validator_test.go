package hw09structvalidator

import (
	"encoding/json"
	"errors"
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
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "Valid User",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Kholod Mariia",
				Age:    26,
				Email:  "kek@gmail.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			name: "Invalid User",
			in: User{
				ID:     "123",
				Name:   "Kholod Mariia",
				Age:    15,
				Email:  "error-email",
				Role:   "user",
				Phones: []string{"123"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: errors.New("length must be 36")},
				{Field: "Age", Err: errors.New("value must be at least 18")},
				{Field: "Email", Err: errors.New("value does not match regexp ^\\w+@\\w+\\.\\w+$")},
				{Field: "Role", Err: errors.New("value must be one of admin,stuff")},
				{Field: "Phones", Err: errors.New("element 0: length must be 11")},
			},
		},
		{
			name: "Valid App",
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			name: "Invalid App",
			in: App{
				Version: "1.0",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: errors.New("length must be 5")},
			},
		},
		{
			name: "Valid Response",
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			name: "Invalid Response",
			in: Response{
				Code: 201,
				Body: "Created",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: errors.New("value must be one of 200,404,500")},
			},
		},
		{
			name:        "Non-struct type",
			in:          "not a struct",
			expectedErr: ErrNotStruct,
		},
		{
			name:        "Struct without validation tags",
			in:          Token{},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)

			if tt.expectedErr == nil && err != nil {
				require.NoError(t, err, "Expected no error")
			} else if tt.expectedErr != nil && err == nil {
				require.Error(t, err, "Expected an error")
			}

			var validationErrors ValidationErrors
			if errors.As(err, &validationErrors) {
				expectedErrors := tt.expectedErr.(ValidationErrors)
				require.Equal(t, len(expectedErrors), len(validationErrors), "Number of errors mismatch")

				for i, expectedErr := range expectedErrors {
					require.Equal(t, expectedErr.Field, validationErrors[i].Field, "Field mismatch")
					require.Equal(t, expectedErr.Err.Error(), validationErrors[i].Err.Error(), "Error message mismatch")
				}
			} else {
				require.ErrorIs(t, err, tt.expectedErr, "Error type mismatch")
			}
		})
	}
}
