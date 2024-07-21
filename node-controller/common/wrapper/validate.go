package wrapper

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var validate = validator.New()

func validateReqBody(reqBody interface{}, expectedTag string) (err error) {
	if reqBody == nil || expectedTag == "" {
		return
	}

	if err = checkByValidator(reqBody, expectedTag); err != nil {
		return
	}

	if checker, ok := reqBody.(ReqBodyValidator); ok {
		if err = checker.CustomCheck(); err != nil {
			return
		}
	}

	return
}

func checkByValidator(reqBody interface{}, expectedTag string) (err error) {
	if err = validate.Struct(reqBody); err == nil {
		return
	}
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return
	}

	var fieldErr = make(map[string]string)
	var reqObjType = reflect.TypeOf(reqBody)
	for _, validationErr := range err.(validator.ValidationErrors) {
		field, ok := reqObjType.Elem().FieldByName(validationErr.StructField())
		if ok {
			fieldErr[field.Tag.Get(expectedTag)] = fmt.Sprintf("%s", validationErr)
		}
	}
	if len(fieldErr) > 0 {
		filedErrDesc, _ := json.Marshal(fieldErr)
		return errors.New(string(filedErrDesc))
	}

	return
}

// ReqBodyValidator 作为validator的补充，用于实现validator无法做到的复杂逻辑校验
type ReqBodyValidator interface {
	CustomCheck() error
}
