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
	ErrNotAStruct                = errors.New("not a struct")
	ErrFieldValidation           = errors.New("validation field error")
	ErrRuleValidation            = errors.New("error in validation rule")
	ErrUnsupportedValidationRule = errors.New("unsupported validation rule")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder

	for _, val := range v {
		sb.WriteString("field '" + val.Field + "': " + val.Err.Error() + "; ")
	}
	return sb.String()
}

func (v *ValidationErrors) Add(field string, err error) {
	*v = append(*v, ValidationError{field, err})
}

type structField struct {
	Type          reflect.Kind
	Name          string
	Value         reflect.Value
	ValidationTag string
	isPublic      bool
}

func Validate(v interface{}) error {
	objValueOf := reflect.ValueOf(v)
	objTypeName := objValueOf.Type().Kind()
	isStruct := objTypeName == reflect.Struct

	if !isStruct {
		return ErrNotAStruct
	}

	validationErrors := ValidationErrors{}

	for field := range iterateOverStructFields(objValueOf) {
		if !field.isPublic || field.ValidationTag == "" {
			continue
		}

		rules := getValidationRules(field.ValidationTag)

		for ruleName, ruleValue := range rules {
			if err := validateSimpleType(field.Type, field.Value, ruleName, ruleValue); err != nil {
				if !errors.Is(err, ErrFieldValidation) {
					return err
				}
				validationErrors.Add(field.Name, err)
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func iterateOverStructFields(objValueOf reflect.Value) <-chan structField {
	ch := make(chan structField)
	go func() {
		defer close(ch)

		valueOf := objValueOf
		fieldCount := valueOf.NumField()

		for i := 0; i < fieldCount; i++ {
			field := valueOf.Type().Field(i)
			fieldKind := field.Type.Kind()
			fieldName := field.Name
			isPublic := field.IsExported()
			validationTag := field.Tag.Get("validate")

			//nolint
			switch fieldKind {
			case reflect.Int, reflect.String:
				val := valueOf.Field(i)
				ch <- structField{fieldKind, fieldName, val, validationTag, isPublic}
			case reflect.Array, reflect.Slice:
				sliceLength := valueOf.Field(i).Len()
				if sliceLength == 0 {
					continue
				}

				sliceKind := valueOf.Field(i).Index(0).Type().Kind()
				for j := 0; j < sliceLength; j++ {
					switch sliceKind {
					case reflect.Int:
					case reflect.String:
						ch <- structField{sliceKind, fieldName, valueOf.Field(i).Index(j), validationTag, isPublic}
					default:
					}
				}
			default:
			}
		}
	}()

	return ch
}

func getValidationRules(tag string) map[string]string {
	rules := make(map[string]string)

	splitedTagByAnd := strings.Split(tag, "|")
	for _, item := range splitedTagByAnd {
		splitedItem := strings.Split(item, ":")
		rules[splitedItem[0]] = splitedItem[1]
	}

	return rules
}

func validateSimpleType(fieldKind reflect.Kind, value reflect.Value, ruleName string, ruleValue string) error {
	//nolint
	switch fieldKind {
	case reflect.Int:
		val := value.Int()
		if err := validateIntValue(int(val), ruleName, ruleValue); err != nil {
			return err
		}
	case reflect.String:
		val := value.String()
		if err := validateStringValue(val, ruleName, ruleValue); err != nil {
			return err
		}
	default:
	}

	return nil
}

func validateIntMin(value int, minValueStr string) error {
	minValue, err := strconv.Atoi(minValueStr)
	if err != nil {
		return ErrRuleValidation
	}

	if value < minValue {
		errMsg := fmt.Sprintf("%d less than min %d", value, minValue)
		return fmt.Errorf("%w - %s", ErrFieldValidation, errMsg)
	}

	return nil
}

func validateIntMax(value int, maxValueStr string) error {
	maxValue, err := strconv.Atoi(maxValueStr)
	if err != nil {
		return ErrRuleValidation
	}

	if value > maxValue {
		errMsg := fmt.Sprintf("%d more than max %d", value, maxValue)
		return fmt.Errorf("%w - %s", ErrFieldValidation, errMsg)
	}

	return nil
}

func validateIntIn(value int, inValueStr string) error {
	splited := strings.Split(inValueStr, ",")
	if len(splited) == 0 {
		errMsg := fmt.Sprintf("cannot parse %s", inValueStr)
		return fmt.Errorf("%w - %s", ErrRuleValidation, errMsg)
	}

	ints := make([]int, len(splited))
	for i, s := range splited {
		intVal, err := strconv.Atoi(s)
		if err != nil {
			return ErrRuleValidation
		}
		ints[i] = intVal
	}

	if !containsIntInSlice(ints, value) {
		errMsg := fmt.Sprintf("%d not contains in tag 'in'", value)
		return fmt.Errorf("%w - %s", ErrFieldValidation, errMsg)
	}

	return nil
}

func containsStringInSlice(s []string, val string) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}

	return false
}

func containsIntInSlice(s []int, val int) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}

	return false
}

func validateStringLen(value string, lenValueStr string) error {
	lenValue, err := strconv.Atoi(lenValueStr)
	if err != nil {
		return ErrRuleValidation
	}

	if len(value) != lenValue {
		errMsg := fmt.Sprintf("length of %s not equal %d", value, lenValue)
		return fmt.Errorf("%w - %s", ErrFieldValidation, errMsg)
	}

	return nil
}

func validateStringRegexp(value string, regexpStr string) error {
	regexp, err := regexp.Compile(regexpStr)
	if err != nil {
		return ErrRuleValidation
	}

	if !regexp.MatchString(value) {
		errMsg := fmt.Sprintf("not match %s", value)
		return fmt.Errorf("%w - %s", ErrFieldValidation, errMsg)
	}

	return nil
}

func validateStringIn(value string, inValueStr string) error {
	splited := strings.Split(inValueStr, ",")

	if !containsStringInSlice(splited, value) {
		errMsg := fmt.Sprintf("%s not contains in tag 'in'", value)
		return fmt.Errorf("%w - %s", ErrFieldValidation, errMsg)
	}

	return nil
}

func validateStringValue(val string, ruleName string, ruleValue string) error {
	var err error

	switch ruleName {
	case "len":
		err = validateStringLen(val, ruleValue)
	case "regexp":
		err = validateStringRegexp(val, ruleValue)
	case "in":
		err = validateStringIn(val, ruleValue)
	default:
		err = fmt.Errorf("%w - %s", ErrUnsupportedValidationRule, ruleName)
	}

	return err
}

func validateIntValue(val int, ruleName string, ruleValue string) error {
	var err error

	switch ruleName {
	case "min":
		err = validateIntMin(val, ruleValue)
	case "max":
		err = validateIntMax(val, ruleValue)
	case "in":
		err = validateIntIn(val, ruleValue)
	default:
		err = fmt.Errorf("%w - %s", ErrUnsupportedValidationRule, ruleName)
	}

	return err
}
