package error

import (
	"fmt"

	"SupCaller/common/locale"
	"SupCaller/common/response"

	"github.com/gin-gonic/gin"
)

// Error 全局错误处理
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Lang    string    `json:"-"`
	Args    []interface{}
}

// New 创建错误对象（使用错误码）
func New(code ErrorCode) *Error {
	return &Error{
		Code: code,
	}
}

// NewWithMessage 创建错误对象（自定义消息）
func NewWithMessage(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewWithArgs 创建错误对象（带参数）
func NewWithArgs(code ErrorCode, args ...interface{}) *Error {
	return &Error{
		Code: code,
		Args: args,
	}
}

// NewWithLang 创建错误对象（指定语言）
func NewWithLang(code ErrorCode, lang string) *Error {
	return &Error{
		Code: code,
		Lang: lang,
	}
}

// NewWithLangAndArgs 创建错误对象（指定语言和参数）
func NewWithLangAndArgs(code ErrorCode, lang string, args ...interface{}) *Error {
	return &Error{
		Code: code,
		Lang: lang,
		Args: args,
	}
}

// GetMessage 获取错误消息（带国际化）
func (e *Error) GetMessage(lang string) string {
	// 如果已经设置了自定义消息，直接返回
	if e.Message != "" {
		if len(e.Args) > 0 {
			return fmt.Sprintf(e.Message, e.Args...)
		}
		return e.Message
	}

	// 根据语言获取国际化消息
	messageKey := e.Code.GetMessageKey()
	return locale.TranslateWithArgs(lang, messageKey, e.Args...)
}

// Handle 全局错误处理函数
func Handle(c *gin.Context, err *Error) {
	// 从上下文获取语言
	lang := locale.GetLangFromContext(c)

	// 如果错误对象中指定了语言，使用指定语言
	if err.Lang != "" {
		lang = err.Lang
	}

	// 获取国际化消息
	message := err.GetMessage(lang)

	// 返回统一格式的错误响应
	c.JSON(200, response.Error(int(err.Code), message))
}

// HandleWithLang 处理错误（指定语言）
func HandleWithLang(c *gin.Context, err *Error, lang string) {
	message := err.GetMessage(lang)
	c.JSON(200, response.Error(int(err.Code), message))
}

// Success 创建成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, response.Success(data))
}

// SuccessWithMessage 创建成功响应（带消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, &response.Response{
		Code:    int(SuccessCode),
		Message: message,
		Data:    data,
	})
}

// AbortWithError 中止请求并返回错误
func AbortWithError(c *gin.Context, err *Error) {
	c.Abort()
	Handle(c, err)
}

// MustNotError 如果有错误则中止请求
func MustNotError(c *gin.Context, err error) {
	if err != nil {
		c.Abort()
		Handle(c, New(InternalServerError))
	}
}

// MustNotErrorWithCode 如果有错误则中止请求（使用指定错误码）
func MustNotErrorWithCode(c *gin.Context, err error, code ErrorCode) {
	if err != nil {
		c.Abort()
		Handle(c, New(code))
	}
}
