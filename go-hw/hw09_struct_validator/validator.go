package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var listOfRules = []string{"len", "in", "max", "min", "regexp", "nested"}

type (
	ValidationError struct {
		Field string
		Err   error
	}

	ValidationErrors []ValidationError
)

var (
	flagValidation      bool
	ErrInputIsNil       = errors.New("error. Input is nil")
	ErrInputIsNotStruct = errors.New("error. Input is not struct")
	ErrValidationString = errors.New("error. Validation input is invalid")
	ErrMismatchedType   = errors.New("error. Type is not supported")
)

const (
	IntArray    = "[]int"
	StringArray = "[]string"
)

func (v ValidationErrors) Error() string {
	res := ""
	for _, val := range v {
		res += fmt.Sprintln(val.Field + ": " + val.Err.Error())
	}
	return res
}

func Validate(v interface{}) (ValidationErrors, error) {
	if v == nil {
		return nil, ErrInputIsNil
	}
	vt := reflect.TypeOf(v).Elem()
	vv := reflect.ValueOf(v).Elem()
	// проверка на то, что входной interface{} - структура.
	if vt.Kind() != reflect.Struct {
		return nil, ErrInputIsNotStruct
	}
	validationErrors := make(ValidationErrors, 0)
	for i := 0; i < vt.NumField(); i++ {
		field := vt.Field(i)
		validate, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}
		// если тег валидации не пустой - не игнорируем
		if validate != "" {
			fieldValue := vv.Field(i)
			rules := strings.Split(validate, "|")
			for j := range rules {
				keyAndValue := strings.Split(rules[j], ":")
				key := keyAndValue[0]
				// проверка, что такой ключ существует
				if !inRule(listOfRules, key) {
					return nil, ErrValidationString
				}
				err := ParseRules(key, field, fieldValue, keyAndValue, &validationErrors)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return validationErrors, nil
}

func ParseRules(
	key string,
	field reflect.StructField,
	fieldValue reflect.Value,
	keyAndValue []string,
	validationErrors *ValidationErrors,
) error {
	// кейсы с валидацией
	switch key {
	case "min":
		{
			err := MinMax(field, fieldValue, keyAndValue, validationErrors, true)
			if err != nil {
				return err
			}
		}
	case "max":
		{
			err := MinMax(field, fieldValue, keyAndValue, validationErrors, false)
			if err != nil {
				return err
			}
		}
	case "in":
		{
			err := In(field, fieldValue, keyAndValue, validationErrors)
			if err != nil {
				return err
			}
		}
	case "len":
		{
			err := Len(field, fieldValue, keyAndValue, validationErrors)
			if err != nil {
				return err
			}
		}
	case "regexp":
		{
			err := Regexp(field, fieldValue, keyAndValue, validationErrors)
			if err != nil {
				return err
			}
		}
	case "nested":
		{
			if fieldValue.Addr().CanInterface() {
				nestedErrors, err := Validate(fieldValue.Addr().Interface())
				if err != nil {
					return err
				}
				*validationErrors = append(*validationErrors, nestedErrors...)
			}
			// проверка на тип
			if field.Type.Kind() != reflect.Struct && field.Type.Kind() != reflect.Slice {
				return ErrMismatchedType
			}
		}
	}
	return nil
}

func MinMax(
	field reflect.StructField,
	fieldValue reflect.Value,
	keyAndValue []string,
	validationErrors *ValidationErrors,
	min bool,
) error {
	// проверка на тип
	if field.Type.Kind() != reflect.Int && field.Type.String() != IntArray {
		return ErrMismatchedType
	}
	// получение min значения из правил
	minValue, err := strconv.Atoi(keyAndValue[1])
	if err != nil {
		return err
	}
	if field.Type.Kind() == reflect.Int {
		// получение значения из элемента
		value := int(fieldValue.Int())

		if min {
			flagValidation = minRule(minValue, value)
		} else {
			flagValidation = maxRule(minValue, value)
		}

		if !flagValidation {
			if min {
				*validationErrors = append(*validationErrors, minMessage(value, minValue, field.Name))
			} else {
				*validationErrors = append(*validationErrors, maxMessage(value, minValue, field.Name))
			}
		}
	} else {
		values, ok := fieldValue.Interface().([]int)
		if !ok {
			return ErrMismatchedType
		}
		for idx, value := range values {
			if min {
				flagValidation = minRule(minValue, value)
			} else {
				flagValidation = maxRule(minValue, value)
			}
			if !flagValidation {
				*validationErrors = append(
					*validationErrors,
					minMessage(
						value,
						minValue,
						fmt.Sprintf("%v[%v]", field.Name, idx)))
				if min {
					*validationErrors = append(
						*validationErrors,
						minMessage(
							value,
							minValue,
							fmt.Sprintf("%v[%v]", field.Name, idx)))
				} else {
					*validationErrors = append(
						*validationErrors,
						maxMessage(
							value,
							minValue,
							fmt.Sprintf("%v[%v]", field.Name, idx)))
				}
			}
		}
	}
	return nil
}

func In(
	field reflect.StructField,
	fieldValue reflect.Value,
	keyAndValue []string,
	validationErrors *ValidationErrors,
) error {
	// проверка на тип
	if field.Type.Kind() != reflect.Int &&
		field.Type.String() != IntArray &&
		field.Type.Kind() != reflect.String &&
		field.Type.String() != StringArray {
		return ErrMismatchedType
	}
	// массив возможных вхождений
	values := strings.Split(keyAndValue[1], ",")
	if field.Type.Kind() == reflect.Int || field.Type.Kind() == reflect.String {
		var value string
		if field.Type.Kind() == reflect.Int {
			value = strconv.Itoa(int(fieldValue.Int()))
		} else {
			value = fieldValue.String()
		}
		flagValidation = inRule(values, value)
		if !flagValidation {
			*validationErrors = append(*validationErrors, inMessage(
				value,
				values,
				field.Name))
		}
	} else {
		var _values []string
		if field.Type.String() == IntArray {
			intValues, ok := fieldValue.Interface().([]int)
			if !ok {
				return ErrMismatchedType
			}
			for i := range intValues {
				number := intValues[i]
				text := strconv.Itoa(number)
				_values = append(_values, text)
			}
		} else {
			var ok bool
			_values, ok = fieldValue.Interface().([]string)
			if !ok {
				return ErrMismatchedType
			}
		}
		for idx, value := range _values {
			flagValidation = inRule(values, value)
			if !flagValidation {
				*validationErrors = append(*validationErrors, inMessage(
					value,
					values,
					fmt.Sprintf("%v[%v]", field.Name, idx)))
			}
		}
	}
	return nil
}

func Len(
	field reflect.StructField,
	fieldValue reflect.Value,
	keyAndValue []string,
	validationErrors *ValidationErrors,
) error {
	// проверка на тип
	if field.Type.Kind() != reflect.String && field.Type.String() != StringArray {
		return ErrMismatchedType
	}
	// получение len значения из правил
	length, err := strconv.Atoi(keyAndValue[1])
	if err != nil {
		return err
	}
	if field.Type.Kind() == reflect.String {
		// получение значения из элемента
		value := fieldValue.String()
		flagValidation = lenRule(length, value)
		if !flagValidation {
			*validationErrors = append(*validationErrors, lenMessage(
				value,
				length,
				field.Name))
		}
	} else {
		values, ok := fieldValue.Interface().([]string)
		if !ok {
			return ErrMismatchedType
		}
		for idx, value := range values {
			flagValidation = lenRule(length, value)
			if !flagValidation {
				*validationErrors = append(*validationErrors, lenMessage(
					value,
					length,
					fmt.Sprintf("%v[%v]", field.Name, idx)))
			}
		}
	}
	return nil
}

func Regexp(
	field reflect.StructField,
	fieldValue reflect.Value,
	keyAndValue []string,
	validationErrors *ValidationErrors,
) error {
	// проверка типа
	if field.Type.Kind() != reflect.String && field.Type.String() != StringArray {
		return ErrMismatchedType
	}
	// получение regexp значения из правил
	regExpPattern := keyAndValue[1]

	if field.Type.Kind() == reflect.String {
		value := fieldValue.String()
		matched, err := regexpRule(regExpPattern, value)
		if err != nil {
			return err
		}
		if !matched {
			*validationErrors = append(*validationErrors, regexpMessage(
				value,
				regExpPattern,
				field.Name))
		}
	} else {
		values, ok := fieldValue.Interface().([]string)
		if !ok {
			return ErrMismatchedType
		}
		for idx, value := range values {
			matched, err := regexpRule(regExpPattern, value)
			if err != nil {
				return err
			}
			if !matched {
				*validationErrors = append(*validationErrors, regexpMessage(
					value,
					regExpPattern,
					fmt.Sprintf("%v[%v]", field.Name, idx)))
			}
		}
	}
	return nil
}
