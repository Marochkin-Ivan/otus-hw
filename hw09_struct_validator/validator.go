package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidValidatorFormat = errors.New("invalid validator format")
	ErrInvalidValidatorType   = errors.New("unsupported validator type")
	ErrInvalidRegexp          = errors.New("invalid regexp")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	for _, err := range v {
		builder.WriteString(fmt.Sprintf("Field: %s, Error: %s\n", err.Field, err.Err))
	}
	return builder.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return errors.New("input is not a struct")
	}

	var validationErrors ValidationErrors
	for i := 0; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i)
		validateTag := fieldType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		fieldName := fieldType.Name
		validators := strings.Split(validateTag, "|")
		for _, validator := range validators {
			err := applyValidator(val.Field(i), validator)
			if err != nil {
				if errors.Is(err, ErrInvalidValidatorFormat) ||
					errors.Is(err, ErrInvalidValidatorType) ||
					errors.Is(err, ErrInvalidRegexp) {
					return fmt.Errorf("field %s: %w", fieldName, err)
				}
				validationErrors = append(validationErrors, ValidationError{Field: fieldName, Err: err})
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func applyValidator(field reflect.Value, validator string) error {
	parts := strings.Split(validator, ":")
	if len(parts) != 2 {
		return ErrInvalidValidatorFormat
	}
	validatorType, validatorValue := parts[0], parts[1]

	switch field.Kind() {
	case reflect.Int:
		return validateInt(field.Int(), validatorType, validatorValue)

	case reflect.String:
		return validateString(field.String(), validatorType, validatorValue)

	case reflect.Slice:
		return validateSlice(field, validator)

	default:
		return ErrInvalidValidatorType
	}
}

func validateInt(value int64, validatorType, validatorValue string) error {
	switch validatorType {
	case "min":
		min, err := strconv.ParseInt(validatorValue, 10, 64)
		if err != nil {
			return ErrInvalidValidatorFormat
		}

		if value < min {
			return fmt.Errorf("value must be at least %d", min)
		}

	case "max":
		max, err := strconv.ParseInt(validatorValue, 10, 64)
		if err != nil {
			return ErrInvalidValidatorFormat
		}

		if value > max {
			return fmt.Errorf("value must be at most %d", max)
		}

	case "in":
		allowedValues := strings.Split(validatorValue, ",")
		for _, v := range allowedValues {
			allowedValue, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return ErrInvalidValidatorFormat
			}

			if value == allowedValue {
				return nil
			}
		}

		return fmt.Errorf("value must be one of %s", allowedValues)

	default:
		return ErrInvalidValidatorType
	}

	return nil
}

func validateString(value, validatorType, validatorValue string) error {
	switch validatorType {
	case "len":
		expectedLen, err := strconv.Atoi(validatorValue)
		if err != nil {
			return ErrInvalidValidatorFormat
		}

		if len(value) != expectedLen {
			return fmt.Errorf("length must be %d", expectedLen)
		}

	case "regexp":
		re, err := regexp.Compile(validatorValue)
		if err != nil {
			return ErrInvalidRegexp
		}

		if !re.MatchString(value) {
			return fmt.Errorf("value does not match regexp: %s", validatorValue)
		}

	case "in":
		allowedValues := strings.Split(validatorValue, ",")
		for _, v := range allowedValues {
			if value == v {
				return nil
			}
		}

		return fmt.Errorf("value must be one of %s", allowedValues)

	default:
		return ErrInvalidValidatorType
	}

	return nil
}

func validateSlice(field reflect.Value, validator string) error {
	for i := 0; i < field.Len(); i++ {
		elem := field.Index(i)
		err := applyValidator(elem, validator)
		if err != nil {
			return err
		}
	}

	return nil
}
