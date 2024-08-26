package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	ErrNotStruct          = errors.New("not a struct")
	ErrInvalidValidateTag = errors.New("invalid validate tag")
	ErrUnsupportedType    = errors.New("unsupported type")
)

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	errorMessages := make([]string, 0, 10)
	for _, err := range v {
		errorMessages = append(errorMessages, fmt.Sprintf("Field: %s, Error: %s", err.Field, err.Err))
	}

	return strings.Join(errorMessages, "\n")
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	var validationErrors ValidationErrors

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if !field.IsExported() {
			continue
		}

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		fieldValue := val.Field(i)
		fieldErrors := validateField(field.Name, fieldValue, tag)
		validationErrors = append(validationErrors, fieldErrors...)
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateField(fieldName string, fieldValue reflect.Value, tag string) ValidationErrors {
	var errors ValidationErrors

	validators := strings.Split(tag, "|")
	for _, validator := range validators {
		parts := strings.SplitN(validator, ":", 2)
		if len(parts) != 2 {
			errors = append(errors, ValidationError{Field: fieldName, Err: ErrInvalidValidateTag})
			continue
		}

		validatorName := parts[0]
		validatorValue := parts[1]

		var err error
		//nolint:exhaustive
		switch fieldValue.Kind() {
		case reflect.Int:
			err = validateInt(fieldValue.Int(), validatorName, validatorValue)
		case reflect.String:
			err = validateString(fieldValue.String(), validatorName, validatorValue)
		case reflect.Slice:
			err = validateSlice(fieldValue, validatorName, validatorValue)
		default:
			err = ErrUnsupportedType
		}

		if err != nil {
			errors = append(errors, ValidationError{Field: fieldName, Err: err})
		}
	}

	return errors
}

func validateInt(value int64, validatorName, validatorValue string) error {
	switch validatorName {
	case "min":
		min, err := strconv.ParseInt(validatorValue, 10, 64)
		if err != nil {
			return ErrInvalidValidateTag
		}
		if value < min {
			return fmt.Errorf("value must be at least %d", min)
		}
	case "max":
		max, err := strconv.ParseInt(validatorValue, 10, 64)
		if err != nil {
			return ErrInvalidValidateTag
		}
		if value > max {
			return fmt.Errorf("value must be at most %d", max)
		}
	case "in":
		values := strings.Split(validatorValue, ",")
		for _, v := range values {
			if strconv.FormatInt(value, 10) == v {
				return nil
			}
		}
		return fmt.Errorf("value must be one of %s", validatorValue)
	default:
		return ErrInvalidValidateTag
	}
	return nil
}

func validateString(value string, validatorName, validatorValue string) error {
	switch validatorName {
	case "len":
		length, err := strconv.Atoi(validatorValue)
		if err != nil {
			return ErrInvalidValidateTag
		}
		if len(value) != length {
			return fmt.Errorf("length must be %d", length)
		}
	case "regexp":
		re, err := regexp.Compile(validatorValue)
		if err != nil {
			return ErrInvalidValidateTag
		}
		if !re.MatchString(value) {
			return fmt.Errorf("value does not match regexp %s", validatorValue)
		}
	case "in":
		values := strings.Split(validatorValue, ",")
		for _, v := range values {
			if value == v {
				return nil
			}
		}
		return fmt.Errorf("value must be one of %s", validatorValue)
	default:
		return ErrInvalidValidateTag
	}
	return nil
}

func validateSlice(value reflect.Value, validatorName, validatorValue string) error {
	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)
		var err error
		//nolint:exhaustive
		switch elem.Kind() {
		case reflect.Int:
			err = validateInt(elem.Int(), validatorName, validatorValue)
		case reflect.String:
			err = validateString(elem.String(), validatorName, validatorValue)
		default:
			return ErrUnsupportedType
		}
		if err != nil {
			return fmt.Errorf("element %d: %w", i, err)
		}
	}
	return nil
}
