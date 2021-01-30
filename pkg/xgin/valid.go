package xgin

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var _validate *validator.Validate

func init() {
	_validate = validator.New()
}

func ContextBindWithValid(ctx *gin.Context, obj interface{}) (err error) {
	err = ctx.ShouldBind(obj)
	if err != nil {
		return
	}

	if _validate != nil {
		err = _validate.Struct(obj)
	}
	return
}

func ContextBindQueryWithValid(ctx *gin.Context, obj interface{}) (err error) {
	err = ctx.ShouldBindQuery(obj)
	if err != nil {
		return
	}

	if _validate != nil {
		err = _validate.Struct(obj)
	}
	return
}

func TelephoneValid(phone string) bool {
	//reg := `^1([387][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	reg := `^1\d{10}$`

	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}
