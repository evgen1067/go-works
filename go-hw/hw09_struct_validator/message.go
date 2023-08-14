package hw09structvalidator

import (
	"errors"
	"fmt"
)

func minMessage(value, minValue int, fieldName string) ValidationError {
	errText := fmt.Sprintf("the number `%v` is less than %v", value, minValue)
	return ValidationError{
		Field: fieldName,
		Err:   errors.New(errText),
	}
}

func maxMessage(value, maxValue int, fieldName string) ValidationError {
	errText := fmt.Sprintf("the number `%v` is greater than %v", value, maxValue)
	return ValidationError{
		Field: fieldName,
		Err:   errors.New(errText),
	}
}

func inMessage(value string, values []string, fieldName string) ValidationError {
	errText := fmt.Sprintf("the value `%v` is not an element of the %v list", value, values)
	return ValidationError{
		Field: fieldName,
		Err:   errors.New(errText),
	}
}

func lenMessage(value string, length int, fieldName string) ValidationError {
	errText := fmt.Sprintf("the length of the `%v` value is different from %v", value, length)
	return ValidationError{
		Field: fieldName,
		Err:   errors.New(errText),
	}
}

func regexpMessage(value, regExpPattern, fieldName string) ValidationError {
	errText := fmt.Sprintf("The value of the expression `%v` does not match the regular expression %v",
		value,
		regExpPattern)
	return ValidationError{
		Field: fieldName,
		Err:   errors.New(errText),
	}
}
