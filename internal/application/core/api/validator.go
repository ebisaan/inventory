package api

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	pg_validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/ebisaan/inventory/internal/application/core/domain"
)

const (
	MapMessageType messageType = iota
	StringMessageType
)

type messageType int

type validate struct {
	validate *pg_validator.Validate
	trans    ut.Translator
}

func newValidate(tagName string) (*validate, error) {
	v := pg_validator.New()
	english := en.New()
	un := ut.New(english, english)
	trans, ok := un.GetTranslator("en")
	if !ok {
		return nil, errors.New("get english translation")
	}
	err := en_translations.RegisterDefaultTranslations(v, trans)
	if err != nil {
		return nil, fmt.Errorf("register default translation(english): %w", err)
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get(tagName), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &validate{
		validate: v,
		trans:    trans,
	}, nil
}

func (v *validate) ValidateStruct(s any) error {
	err := v.validate.Struct(s)
	if err != nil {
		returnErr := domain.ValidationError{}
		var ve pg_validator.ValidationErrors
		if errors.As(err, &ve) {
			returnErr.ValidationErrorTranslations = ve.Translate(v.trans)
			return returnErr
		} else {
			return err
		}
	}

	return nil
}
