package hw09structvalidator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	UserRole string

	User struct {
		ID     string   `validate:"len:12"`
		Name   string   `validate:"len:6"`
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		Skill  Level    `validate:"nested"`
	}

	Level struct {
		ID   string `validate:"len:12"`
		Name string `validate:"in:senior,middle,junior"`
	}

	FailStruct struct {
		Code int `validate:"error:200,404,500"`
		Body string
	}
)

func TestFailValidate(t *testing.T) {
	str, num, fail := "user", 1067, FailStruct{
		Code: 404,
		Body: "body",
	}
	tests := []struct {
		input           interface{}
		expected        error
		expectedMessage string
	}{
		{
			input:           nil,
			expected:        ErrInputIsNil,
			expectedMessage: "",
		},
		{
			input:           &str,
			expected:        ErrInputIsNotStruct,
			expectedMessage: "",
		},
		{
			input:           &num,
			expected:        ErrInputIsNotStruct,
			expectedMessage: "",
		},
		{
			input:           &fail,
			expected:        ErrValidationString,
			expectedMessage: "",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			result, err := Validate(tt.input)
			require.ErrorIs(t, err, tt.expected)
			if tt.expectedMessage != "" {
				require.Equal(t, tt.expectedMessage, result.Error())
			}
		})
	}
}

var user = User{
	ID:     "123456789012",
	Name:   "sixsix",
	Age:    21,
	Email:  "test@mail.ru",
	Role:   "admin",
	Phones: []string{"89912192614"},
	Skill: Level{
		ID:   "121212121212",
		Name: "middle",
	},
}

func TestLength(t *testing.T) {
	tests := []struct {
		id       string
		name     string
		expected string
	}{
		{
			id:       "1",
			name:     "sixsix",
			expected: "ID: the length of the `1` value is different from 12\n",
		},
		{
			id:       "123456789012",
			name:     "six",
			expected: "Name: the length of the `six` value is different from 6\n",
		},
		{
			id:       "123456789012",
			name:     "sixsix",
			expected: "",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			user.ID = tt.id
			user.Name = tt.name
			result, err := Validate(&user)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result.Error())
		})
	}
}

func TestRegex(t *testing.T) {
	tests := []struct {
		email    string
		expected string
	}{
		{
			email: "testmail.ru",
			expected: "Email: The value of the expression `testmail.ru`" +
				" does not match the regular expression ^\\w+@\\w+\\.\\w+$\n",
		},
		{
			email:    "test@mail.ru",
			expected: "",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			user.Email = tt.email
			result, err := Validate(&user)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result.Error())
		})
	}
}

func TestInArray(t *testing.T) {
	tests := []struct {
		role     UserRole
		expected string
	}{
		{
			role:     "admin",
			expected: "",
		},
		{
			role:     "user",
			expected: "Role: the value `user` is not an element of the [admin stuff] list\n",
		},
		{
			role:     "stuff",
			expected: "",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			user.Role = tt.role
			result, err := Validate(&user)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result.Error())
		})
	}
}

func TestMinMax(t *testing.T) {
	tests := []struct {
		age      int
		expected string
	}{
		{
			age:      12,
			expected: "Age: the number `12` is less than 18\n",
		},
		{
			age:      21,
			expected: "",
		},
		{
			age:      60,
			expected: "Age: the number `60` is greater than 50\n",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			user.Age = tt.age
			result, err := Validate(&user)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result.Error())
		})
	}
}
