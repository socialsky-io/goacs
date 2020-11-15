package request

import (
	"errors"
	"github.com/gin-gonic/gin"
	locen "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"log"
	"reflect"
	"strings"
)

type ApiValidator struct {
	Context    *gin.Context
	Data       interface{}
	validator  *validator.Validate
	Errors     map[string]interface{}
	Translator ut.Translator
}

/*
 */
func NewApiValidator(ctx *gin.Context, data interface{}) *ApiValidator {
	validatorInstance := validator.New()
	// From docs
	validatorInstance.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	uTranslator, _ := ut.New(locen.New()).GetTranslator("en")
	en.RegisterDefaultTranslations(validatorInstance, uTranslator)

	return &ApiValidator{
		Context:    ctx,
		Data:       data,
		validator:  validatorInstance,
		Errors:     map[string]interface{}{},
		Translator: uTranslator,
	}
}

func (v *ApiValidator) Validate() error {
	//Clear error
	v.Errors = map[string]interface{}{}

	err := v.validator.Struct(v.Data)

	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Println(err)
			return err
		}

		err := err.(validator.ValidationErrors)
		for _, fieldError := range err {
			log.Println(fieldError.Field()) // by passing alt name to ReportError like below
			log.Println(fieldError.Translate(v.Translator))
			v.Errors[fieldError.Field()] = fieldError.Translate(v.Translator)
		}
		//TO DO RETURN ERROR
		return errors.New("Validation Error")
	}

	return nil
}
