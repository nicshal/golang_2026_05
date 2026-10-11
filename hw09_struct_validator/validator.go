package hw09structvalidator

import (
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

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error in field %s: %s", e.Field, e.Err)
}

func (errs ValidationErrors) Error() string {
	errMsg := make([]string, 0)
	for _, err := range errs {
		errMsg = append(errMsg, err.Error())
	}
	return strings.Join(errMsg, "; ")
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a struct")
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if field.PkgPath != "" {
			continue // Ignore unexported fields
		}

		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue // Ignore fields without validate tag
		}

		fieldValue := val.Field(i)
		if !fieldValue.IsValid() || !fieldValue.CanInterface() {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err := validateInt(fieldValue.Int(), validateTag)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			}
		case reflect.String:
			err := validateString(fieldValue.String(), validateTag)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			}
		case reflect.Slice:
			err := validateSlice(fieldValue.Interface(), validateTag)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err})
			}
		case reflect.Invalid, reflect.Bool, reflect.Uint, reflect.Uint8,
			reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
			reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
			reflect.Pointer, reflect.Struct, reflect.UnsafePointer:
			panic("unhandled default case")
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateInt(value int64, tag string) error {
	validators := strings.Split(tag, "|")
	for _, validator := range validators {
		if strings.HasPrefix(validator, "min:") {
			minV, err := strconv.Atoi(strings.TrimPrefix(validator, "min:"))
			if err != nil {
				return err
			}
			if value < int64(minV) {
				return fmt.Errorf("value %d is less than minimum %d", value, minV)
			}
		}
		if strings.HasPrefix(validator, "max:") {
			maxV, err := strconv.Atoi(strings.TrimPrefix(validator, "max:"))
			if err != nil {
				return err
			}
			if value > int64(maxV) {
				return fmt.Errorf("value %d is greater than maximum %d", value, maxV)
			}
		}
		if strings.HasPrefix(validator, "in:") {
			allowedValues := strings.Split(strings.TrimPrefix(validator, "in:"), ",")
			for _, allowed := range allowedValues {
				if strconv.Itoa(int(value)) == allowed {
					return nil
				}
			}
			return fmt.Errorf("value %d is not in allowed set %v", value, allowedValues)
		}
	}

	return nil
}

func validateString(value string, tag string) error {
	validators := strings.Split(tag, "|")
	for _, validator := range validators {
		if strings.HasPrefix(validator, "len:") {
			length, err := strconv.Atoi(strings.TrimPrefix(validator, "len:"))
			if err != nil {
				return err
			}
			if len(value) != length {
				return fmt.Errorf("value %s does not have length %d", value, length)
			}
		}
		if strings.HasPrefix(validator, "regexp:") {
			re, err := regexp.Compile(strings.TrimPrefix(validator, "regexp:"))
			if err != nil {
				return err
			}
			if !re.MatchString(value) {
				return fmt.Errorf("value %s does not match regexp %s", value, re.String())
			}
		}
		if strings.HasPrefix(validator, "in:") {
			allowedValues := strings.Split(strings.TrimPrefix(validator, "in:"), ",")
			for _, allowed := range allowedValues {
				if value == allowed {
					return nil
				}
			}
			return fmt.Errorf("value %s is not in allowed set %v", value, allowedValues)
		}
	}

	return nil
}

func validateSlice(value interface{}, tag string) error {
	sliceValue := reflect.ValueOf(value)
	for i := 0; i < sliceValue.Len(); i++ {
		elem := sliceValue.Index(i)
		switch elem.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err := validateInt(elem.Int(), tag)
			if err != nil {
				return err
			}
		case reflect.String:
			err := validateString(elem.String(), tag)
			if err != nil {
				return err
			}
		case reflect.Invalid, reflect.Bool, reflect.Uint, reflect.Uint8, reflect.Slice,
			reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
			reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
			reflect.Pointer, reflect.Struct, reflect.UnsafePointer:
			panic("unhandled default case")
		}
	}

	return nil
}
