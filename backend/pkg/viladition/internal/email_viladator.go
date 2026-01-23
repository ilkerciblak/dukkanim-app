package internal

import (
	"regexp"
	"strings"
)

type EmailViladatorInterface interface {
	GetErrors() []string
	GetValue() string
	IsValid() bool
	Required() EmailViladatorInterface
	Validate() EmailViladatorInterface
	validatePattern() EmailViladatorInterface
	validateLength() EmailViladatorInterface
	// validateDomain() EmailViladatorInterface
}

type EmailViladator struct {
	errors        []string
	value         string
	pattern       string
	local_address string
	domain_name   string
}

func NewEmailViladator(val any) *EmailViladator {
	value, k := val.(string)
	if !k {
		return &EmailViladator{
			errors:  []string{"Input should be in string type"},
			pattern: `^([a-zA-Z0-9]+(?:[._-][a-zA-Z0-9]+)*)@([a-zA-Z0-9]+(?:[.-][a-zA-Z0-9]+)*)(\.[a-zA-Z]{2,}(?:\.[a-zA-Z]{2,})*)$`,
		}
	}
	return &EmailViladator{
		value:   value,
		errors:  make([]string, 0),
		pattern: `^([a-zA-Z0-9]+(?:[._-][a-zA-Z0-9]+)*)@([a-zA-Z0-9]+(?:[.-][a-zA-Z0-9]+)*)(\.[a-zA-Z]{2,}(?:\.[a-zA-Z]{2,})*)$`,
	}
}

func (v EmailViladator) GetErrors() []string {
	return v.errors
}

func (v EmailViladator) GetValue() string {
	return v.value
}

func (v EmailViladator) IsValid() bool {
	return len(v.errors) == 0
}

func (v EmailViladator) Required() EmailViladatorInterface {
	if strings.TrimSpace(v.value) == "" {
		v.errors = append(v.errors, "Email is required")
	}

	return &v
}

func (v *EmailViladator) Validate() EmailViladatorInterface {
	v.validatePattern().
		validateLength()
	return v
}

func (v *EmailViladator) validatePattern() EmailViladatorInterface {
	r, err := regexp.Compile(v.pattern)
	if err != nil {
		v.errors = append(v.errors, err.Error())
	}

	if matches := r.MatchString(v.value); !matches {
		v.errors = append(v.errors, "Email is in invalid format")
	}

	mathes := r.FindStringSubmatch(v.value)
	if len(mathes) != 0 {
		v.local_address = mathes[1]
		v.domain_name = mathes[2]
	}

	return v
}

func (v *EmailViladator) validateLength() EmailViladatorInterface {
	if len(v.local_address) > 64 {
		v.errors = append(v.errors, "Email username exceeding max limit")
	}

	if len(v.domain_name) > 254 {
		v.errors = append(v.errors, "Email domain exceeding max limits")
	}

	return v
}

/*
RFC
/^((?!\.)[\w\-_.]*[^.])(@\w+)(\.\w+(\.\w+)?[^.\W])$/gm

https://regex101.com/r/SOgUIV/2
*/
