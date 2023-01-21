package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID      string `json:"id" validate:"len:36"`
		Name    string
		Age     int      `validate:"min:18|max:50"`
		Email   string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role    UserRole `validate:"in:admin,stuff"`
		Phones  []string `validate:"len:11"`
		meta    json.RawMessage
		innerID string `validate:"len:36"`
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

	WrongRuleLenString struct {
		WrongRule string `validate:"len:S"`
	}

	WrongRuleRegexpString struct {
		WrongRule string `validate:"regexp:/\\"`
	}

	WrongRuleMinInt struct {
		WrongRule int `validate:"min:wrong"`
	}

	WrongRuleMaxInt struct {
		WrongRule int `validate:"max:wrong"`
	}

	WrongRuleInInt struct {
		WrongRule int `validate:"in:wrong"`
	}

	UnsupportedRuleInt struct {
		UnsupportedRule int `validate:"unsupported:wrong"`
	}

	UnsupportedRuleString struct {
		UnsupportedRule string `validate:"unsupported:wrong"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			"not a struct",
			"i an not a struct",
			ErrNotAStruct,
		},
		{
			"no validation tag",
			Token{
				Header:    []byte{1, 2, 3, 4},
				Payload:   []byte{4, 3, 2, 1},
				Signature: []byte{},
			},
			nil,
		},
		{
			"validation tag and other tag",
			Response{
				Code: 200,
				Body: "Test",
			},
			nil,
		},
		{
			"error in validation rule min type int",
			WrongRuleMinInt{
				WrongRule: 1,
			},
			ErrRuleValidation,
		},
		{
			"error in validation rule max type int",
			WrongRuleMaxInt{
				WrongRule: 1,
			},
			ErrRuleValidation,
		},
		{
			"error in validation rule in type int",
			WrongRuleInInt{
				WrongRule: 1,
			},
			ErrRuleValidation,
		},
		{
			"error in validation rule len type string",
			WrongRuleLenString{
				WrongRule: "12345",
			},
			ErrRuleValidation,
		},
		{
			"error in validation rule regexp type string",
			WrongRuleRegexpString{
				WrongRule: "123",
			},
			ErrRuleValidation,
		},
		{
			"unsupported validation rule int",
			UnsupportedRuleInt{
				UnsupportedRule: 1,
			},
			fmt.Errorf("%w - %s", ErrUnsupportedValidationRule, "unsupported"),
		},
		{
			"unsupported validation rule string",
			UnsupportedRuleString{
				UnsupportedRule: "123",
			},
			fmt.Errorf("%w - %s", ErrUnsupportedValidationRule, "unsupported"),
		},
		{
			"all fields errors",
			User{
				ID:      "0",
				Name:    "Ivan",
				Age:     1,
				Email:   "email",
				Role:    "user",
				Phones:  []string{"123", "321"},
				meta:    nil,
				innerID: "123",
			},
			ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   fmt.Errorf("%w - length of 0 not equal 36", ErrFieldValidation),
				},
				ValidationError{
					Field: "Age",
					Err:   fmt.Errorf("%w - 1 less than min 18", ErrFieldValidation),
				},
				ValidationError{
					Field: "Email",
					Err:   fmt.Errorf("%w - not match email", ErrFieldValidation),
				},
				ValidationError{
					Field: "Role",
					Err:   fmt.Errorf("%w - user not contains in tag 'in'", ErrFieldValidation),
				},
				ValidationError{
					Field: "Phones",
					Err:   fmt.Errorf("%w - length of 123 not equal 11", ErrFieldValidation),
				},
				ValidationError{
					Field: "Phones",
					Err:   fmt.Errorf("%w - length of 321 not equal 11", ErrFieldValidation),
				},
			},
		},
		{
			"no errors",
			User{
				ID:      "000000000000000000000000000000000001",
				Name:    "Ivan",
				Age:     18,
				Email:   "1@1.ru",
				Role:    "admin",
				Phones:  []string{"12345678910", "10987654321", "10987654321", "10987654321"},
				meta:    nil,
				innerID: "123",
			},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}
