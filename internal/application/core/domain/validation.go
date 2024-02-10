package domain

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	FieldErrorMessages map[string]string
}

func (err ValidationError) Error() string {
	returnMsg := ""
	for k, v := range err.FieldErrorMessages {
		returnMsg += fmt.Sprintf("%s: %s\n", k, v)
	}
	return returnMsg
}

func (err ValidationError) FieldMessages() any {
	validationErrors := map[string]string{}
	for k, v := range err.FieldErrorMessages {
		tags := strings.Split(k, ".")
		if len(tags) > 1 {
			validationErrors[tags[len(tags)-1]] = v
		} else {
			validationErrors[k] = v
		}
	}
	return validationErrors
}
