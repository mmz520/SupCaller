package validate

import (
	"reflect"
	"regexp"
	"strings"
)

// ValidationError 验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validator 数据校验器接口
type Validator interface {
	Validate() []*ValidationError
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

// ValidateMobile 验证手机号格式
func ValidateMobile(mobile string) bool {
	pattern := `^1[3-9]\d{9}$`
	match, _ := regexp.MatchString(pattern, mobile)
	return match
}

// ValidateRequired 验证必填字段
func ValidateRequired(value interface{}) bool {
	if value == nil {
		return false
	}
	
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return strings.TrimSpace(v.String()) != ""
	case reflect.Ptr, reflect.Interface:
		return !v.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() > 0
	default:
		return true
	}
}