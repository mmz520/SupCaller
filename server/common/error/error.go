package error

import (
	"net/http"
	"SupCaller/common/response"
	
	"github.com/gin-gonic/gin"
)

// Error 全局错误处理
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// New 创建错误对象
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Handle 全局错误处理函数
func Handle(c *gin.Context, err *Error) {
	c.JSON(http.StatusOK, response.Error(err.Code, err.Message))
}

