package compact

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Marshal converts a struct to a comma-separated string
func Marshal(v any) (string, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return "", errors.New("input must be a struct")
	}

	var parts []string

	typ := val.Type()
	for i := range val.NumField() {
		field := val.Field(i)
		if !field.CanInterface() {
			continue // Skip unexported fields
		}

		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			parts = append(parts, strconv.FormatInt(field.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			parts = append(parts, strconv.FormatUint(field.Uint(), 10))
		case reflect.Float32, reflect.Float64:
			parts = append(parts, strconv.FormatFloat(field.Float(), 'f', -1, 64))
		case reflect.String:
			parts = append(parts, field.String())
		case reflect.Bool:
			parts = append(parts, strconv.FormatBool(field.Bool()))
		default:
			return "", fmt.Errorf("unsupported field type: %s", typ.Field(i).Name)
		}
	}

	return strings.Join(parts, ","), nil
}

// Unmarshal parses a comma-separated string into a struct
func Unmarshal(data string, v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() || val.Elem().Kind() != reflect.Struct {
		return errors.New("input must be a pointer to a struct")
	}

	val = val.Elem()

	parts := strings.Split(data, ",")
	if len(parts) != val.NumField() {
		return fmt.Errorf("mismatched field count: expected %d, got %d", val.NumField(), len(parts))
	}

	for i := range val.NumField() {
		field := val.Field(i)
		if !field.CanSet() {
			continue // Skip unexported fields
		}

		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(parts[i], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid integer for field %s: %w", val.Type().Field(i).Name, err)
			}

			field.SetInt(n)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(parts[i], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid unsigned integer for field %s: %w", val.Type().Field(i).Name, err)
			}

			field.SetUint(n)
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(parts[i], 64)
			if err != nil {
				return fmt.Errorf("invalid float for field %s: %w", val.Type().Field(i).Name, err)
			}

			field.SetFloat(n)
		case reflect.String:
			field.SetString(parts[i])
		case reflect.Bool:
			b, err := strconv.ParseBool(parts[i])
			if err != nil {
				return fmt.Errorf("invalid bool for field %s: %w", val.Type().Field(i).Name, err)
			}

			field.SetBool(b)
		default:
			return fmt.Errorf("unsupported field type: %s", val.Type().Field(i).Name)
		}
	}

	return nil
}
