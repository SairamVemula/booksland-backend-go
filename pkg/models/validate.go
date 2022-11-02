package models

import (
	"fmt"
	"log"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// ValidationError wraps the validator FieldError so we do not
// expose this to outside code
type ValidationError struct {
	validator.FieldError
}

// Error provides the string format of the validation error
func (v ValidationError) Error() string {
	var err string

	switch v.Tag() {
	case "required":
		err = fmt.Sprintf("%s is required", v.Field())
	case "min":
		err = fmt.Sprintf("%s should be atleast %s charactars", v.Field(), v.Param())
	case "max":
		err = fmt.Sprintf("%s should be atmost %s charactars", v.Field(), v.Param())
	case "email":
		err = fmt.Sprintf("Enter a valid %s", v.Field())
	case "numeric":
		err = fmt.Sprintf("%s should only have numeric", v.Field())
	case "passwd":
		err = fmt.Sprintf("%s should have Minimum eight characters, at least one uppercase letter, one lowercase letter, one number and one special character", v.Field())
	}

	// fmt.Sprintf(
	// 		"key: '%s' Error: Field validation for '%s' failed on the '%s' tag",
	// 		v.Namespace(),
	// 		v.Field(),
	// 		v.Tag())

	return err
}

// ValidationErrors is a wrapper for list of ValidationError
type ValidationErrors []ValidationError

// Errors convert the ValidationErrors slice into string slice
func (v ValidationErrors) Errors() []string {
	errs := []string{}
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

// Validation is the type for validator
type Validation struct {
	validate *validator.Validate
}

// NewValidation returns a Validator instance
func NewValidation() *Validation {
	validate := validator.New()
	err := validate.RegisterValidation("passwd", Passwd)
	if err != nil {
		log.Println(err.Error())
	}
	return &Validation{validate}
}

// Struct method validates the given struct based on the validate tags
// and returns validation error if any
func (v *Validation) Struct(i interface{}) ValidationErrors {
	errs := v.validate.Struct(i)
	if errs == nil {
		return nil
	}

	var returnErrs ValidationErrors
	for _, err := range errs.(validator.ValidationErrors) {
		// cast the FieldError into our ValidationError and append to the slice
		ve := ValidationError{err.(validator.FieldError)}
		returnErrs = append(returnErrs, ve)
	}
	return returnErrs
}
func (v *Validation) StructExcept(i interface{}, fields ...string) ValidationErrors {
	errs := v.validate.StructExcept(i, fields...)
	if errs == nil {
		return nil
	}

	var returnErrs ValidationErrors
	for _, err := range errs.(validator.ValidationErrors) {
		// cast the FieldError into our ValidationError and append to the slice
		ve := ValidationError{err.(validator.FieldError)}
		returnErrs = append(returnErrs, ve)
	}
	return returnErrs
}

var Passwd = func(fl validator.FieldLevel) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	s := fl.Field().String()
	if len(s) >= 8 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

var LocationValidation = func(field validator.FieldLevel) bool {
	inter := field.Field()
	slice, ok := inter.Interface().([]string)
	if !ok {
		return false
	}
	return len(slice) == 2
}
