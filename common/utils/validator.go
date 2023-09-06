package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

type Validator interface {
	Validate(data interface{}) string
}

// ZhValidator 参数校验器
type ZhValidator struct {
	Validator *validator.Validate
	Uni       *ut.UniversalTranslator
	Trans     map[string]ut.Translator
}

// NewZhValidator
// @Description: 获取参数校验器
// @return *ZhValidator
func NewZhValidator() *ZhValidator {
	v := ZhValidator{}
	translator := zh.New()
	v.Uni = ut.New(translator, translator)
	v.Validator = validator.New()
	zhTrans, _ := v.Uni.GetTranslator("translator")
	v.Trans = make(map[string]ut.Translator)
	v.Trans["translator"] = zhTrans

	err := v.Validator.RegisterValidation("password", passwordValidation)
	if err != nil {
		panic(fmt.Sprintf("校验器注册错误：%v", err))
	}
	// 自定义翻译消息
	_ = v.Validator.RegisterTranslation("password", zhTrans, func(ut ut.Translator) error {
		return ut.Add("password", "{0}密码必须以字母开头，长度在6~18之间，只能包含字母、数字和下划线", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})
	err = v.Validator.RegisterValidation("zh_cn_phone", chinesePhoneValidation)
	if err != nil {
		panic(fmt.Sprintf("校验器注册错误：%v", err))
	}
	// 自定义翻译消息
	_ = v.Validator.RegisterTranslation("zh_cn_phone", zhTrans, func(ut ut.Translator) error {
		return ut.Add("zh_cn_phone", "{0}请填写中国手机号", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("zh_cn_phone", fe.Field())
		return t
	})
	err = zhTranslations.RegisterDefaultTranslations(v.Validator, zhTrans)
	if err != nil {
		return nil
	}
	return &v
}

// Validate
// @Description: 中文校验返回错误信息翻译
// @receiver v
// @param data
// @param lang
// @return string
func (v *ZhValidator) Validate(data interface{}) string {
	err := v.Validator.Struct(data)
	if err == nil {
		return ""
	}

	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if ok {
		transData := errs.Translate(v.Trans["zh"])
		s := strings.Builder{}
		for _, v := range transData {
			s.WriteString(v)
			s.WriteString(" ")
		}
		return s.String()
	}

	var invalid *validator.InvalidValidationError
	ok = errors.As(err, &invalid)
	if ok {
		return invalid.Error()
	}

	return ""
}

// passwordValidation
// @Description: 自定义密码校验器 -> 密码必须以字母开头，长度在6~18之间，只能包含字母、数字和下划线
// @param fl
// @return bool
func passwordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// 密码必须以字母开头，长度在6~18之间，只能包含字母、数字和下划线
	pattern := `^[a-zA-Z]\w{5,17}$`
	match, _ := regexp.MatchString(pattern, password)
	return match
}

// chinesePhoneValidation
// @Description: 中国手机号校验
// @param fl
// @return bool
func chinesePhoneValidation(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// 中国手机号格式为11位数字，以1开头
	pattern := `^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`
	match, _ := regexp.MatchString(pattern, phone)
	return match
}
