package internal

import (
	"fmt"

	"regexp"
	"strings"
	"unicode/utf8"
)

type IStringValidator interface {
	GetErrors() []string
	GetFirstErr() string
	GetValue() string
	IsValid() bool
	Required() IStringValidator
	MinLength(val int) IStringValidator
	MaxLength(val int) IStringValidator
	Pattern(pattern string) IStringValidator
}

type StringViladator struct {
	errors []string
	value  string
}

func NewStringViladator(val any) *StringViladator {
	value, k := val.(string)
	if !k {
		return &StringViladator{
			errors: []string{"Input should be in string type"},
		}
	}
	return &StringViladator{
		value:  value,
		errors: make([]string, 0),
	}
}

func (v *StringViladator) GetFirstErr() string {
	if len(v.errors) != 0 {
		return v.errors[0]
	}
	return ""
}

func (v StringViladator) GetErrors() []string {

	return v.errors
}

func (v StringViladator) IsValid() bool {
	return len(v.errors) == 0
}

func (v StringViladator) GetValue() string {
	return v.value
}

func (v *StringViladator) Required() IStringValidator {
	if strings.TrimSpace(v.value) == "" {
		v.errors = append(v.errors, "Field cannot be empty")
	}
	return v
}
func (v *StringViladator) MinLength(val int) IStringValidator {
	if utf8.RuneCountInString(v.value) < val {
		v.errors = append(v.errors, fmt.Sprintf("Minimum length is %d characters", val))
	}
	return v
}
func (v *StringViladator) MaxLength(val int) IStringValidator {
	if utf8.RuneCountInString(v.value) > val {
		v.errors = append(v.errors, fmt.Sprintf("Maximum length is %d characters", val))
	}

	return v
}
func (v *StringViladator) Pattern(pattern string) IStringValidator {

	if match, err := regexp.MatchString(pattern, v.value); err != nil {
		v.errors = append(v.errors, fmt.Sprintf("[Internal Server Error]:%v", err))
	} else if !match {
		v.errors = append(v.errors, fmt.Sprintf("Input does not match the pattern:%v", pattern))
	}

	return v
}
