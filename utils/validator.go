package utils

// 将验证器错误翻译成中文

import (
	"encoding/json"
	"fmt"
	"github.com/derekAHua/goLib/errors"
	"reflect"
	"strconv"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni   *ut.UniversalTranslator
	trans ut.Translator
)

func InitValidator() {
	// 注册翻译器
	z := zh.New()
	uni = ut.New(z, z)

	trans, _ = uni.GetTranslator("zh")

	// 获取gin的校验器
	validate := binding.Validator.Engine().(*validator.Validate)

	// 注册自定义函数
	_ = validate.RegisterValidation("runeLen", MaxRuneLength)
	_ = validate.RegisterValidation("phone", PhoneFormat)

	// 注册tag
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		if len(name) == 0 {
			return fld.Name
		}
		return name
	})
	// 注册翻译器
	_ = zhTranslations.RegisterDefaultTranslations(validate, trans)
}

// Translate 翻译错误信息
func Translate(err error) error {
	resErr := errors.Err{ErrNo: -1, ErrMsg: "参数错误"}

	switch err.(type) {
	case validator.ValidationErrors:
		validationErrors, _ := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			for _, e := range validationErrors {
				resErr.ErrMsg = e.Translate(trans)
				break
			}
		}
	case *json.UnmarshalTypeError:
		unmarshalErr := err.(*json.UnmarshalTypeError)
		resErr.ErrMsg = fmt.Sprintf("%s参数类型错误", unmarshalErr.Field)
	case *json.SyntaxError:
		resErr.ErrMsg = "入参json格式不正确"
	default:
		// 未知类型错误直接返回
		resErr.ErrMsg = err.Error()
	}

	return resErr
}

// MaxRuneLength 校验字段长度
var MaxRuneLength validator.Func = func(fl validator.FieldLevel) bool {
	str, ok := fl.Field().Interface().(string)
	if ok {
		length, err := strconv.Atoi(fl.Param())
		if err != nil {
			return false
		}
		if len([]rune(str)) > length {
			return false
		}
	}
	return true
}

// PhoneFormat 校验字段电话格式
var PhoneFormat validator.Func = func(fl validator.FieldLevel) bool {
	phoneNum, ok := fl.Field().Interface().(int)
	if ok {
		return CheckPhoneFormat(strconv.Itoa(phoneNum))
	}
	str, ok := fl.Field().Interface().(string)
	if ok {
		return CheckPhoneFormat(str)
	}
	return false
}

func CheckPhoneFormat(p string) (ret bool) {
	if len(p) != 11 {
		return
	}

	for _, v := range p {
		if !unicode.IsDigit(v) {
			return
		}
	}

	var phoneList = map[string]struct{}{
		"134": {}, "135": {}, "136": {}, "137": {}, "138": {}, "139": {},
		"147": {}, "148": {}, "150": {}, "151": {}, "152": {}, "157": {}, "158": {},
		"159": {}, "178": {}, "182": {}, "183": {}, "184": {}, "187": {}, "188": {},
		"198": {}, "165": {}, "172": {}, "195": {}, "197": {}, "130": {}, "131": {},
		"132": {}, "145": {}, "146": {}, "155": {}, "156": {}, "166": {}, "167": {},
		"175": {}, "176": {}, "185": {}, "186": {}, "169": {}, "196": {}, "133": {},
		"141": {}, "153": {}, "162": {}, "173": {}, "174": {}, "177": {}, "180": {},
		"181": {}, "189": {}, "191": {}, "199": {}, "190": {}, "192": {}, "193": {},
		"170": {}, "171": {}, "116": {}, "115": {}, "111": {},
	}
	prefix := p[0:3]
	if _, ok := phoneList[prefix]; !ok {
		return
	}

	ret = true

	return
}
