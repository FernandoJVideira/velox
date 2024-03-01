package velox

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

type Validation struct {
	Data   url.Values
	Errors map[string]string
}

func (v *Velox) Validator(data url.Values) *Validation {
	return &Validation{
		Data:   data,
		Errors: make(map[string]string),
	}
}

func (v *Validation) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validation) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validation) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

func (v *Validation) Required(r *http.Request, fields ...string) {
	for _, field := range fields {
		val := r.Form.Get(field)
		if strings.TrimSpace(val) == "" {
			v.AddError(field, "This field is required")
		}
	}
}

func (v *Validation) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validation) IsEmail(field, val string) {
	if !govalidator.IsEmail(val) {
		v.AddError(field, "Invalid email address")
	}
}

func (v *Validation) IsInt(field, val string) {
	_, err := strconv.Atoi(val)
	if err != nil {
		v.AddError(field, "This field must be an integer")
	}
}

func (v *Validation) IsFloat(field, val string) {
	_, err := strconv.ParseFloat(val, 64)
	if err != nil {
		v.AddError(field, "This field must be a floating point number")
	}
}

func (v *Validation) IsDate(field, val string) {
	_, err := time.Parse("02-01-2006", val)
	if err != nil {
		v.AddError(field, "This field must be a date in the format DD-MM-YYYY")
	}
}

func (v *Validation) NoSpaces(field, val string) {
	if govalidator.HasWhitespace(val) {
		v.AddError(field, "This field must not contain any spaces")
	}
}
