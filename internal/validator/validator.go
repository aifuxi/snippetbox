package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	// FieldErrors map 如果为 nil 则需要初始化
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	// 如果key对应的错误不存在则新增
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField 检查field，如果未通过校验，则加到错误中
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt 判断 value 是否在 permittedValues 中
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}

	return false
}

// 邮箱正则 https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// 正则匹配
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
