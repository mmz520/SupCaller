package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"SupCaller/common/locale"
	"SupCaller/common/logger"

	"github.com/gin-gonic/gin"
)

// bodyLogWriter 用于记录请求体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// OperationLogMiddleware 操作日志中间件
func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 获取请求信息
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

		// 获取请求头信息
		headers := make(map[string]string)
		for k, v := range c.Request.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}
		headersJSON, _ := json.Marshal(headers)

		// 获取查询参数
		queryParams := c.Request.URL.Query().Encode()

		// 获取客户端IP
		clientIP := c.ClientIP()

		// 获取用户信息（从上下文）
		var userID uint
		var username string
		if user, exists := c.Get("user"); exists {
			if u, ok := user.(map[string]interface{}); ok {
				if id, ok := u["id"].(uint); ok {
					userID = id
				}
				if name, ok := u["username"].(string); ok {
					username = name
				}
			}
		}

		// 获取语言和时区
		lang := locale.GetLangFromContext(c)
		timezone := c.GetHeader(locale.TimezoneHeader)
		if timezone == "" {
			timezone = locale.DefaultTimezone
		}

		// 包装响应writer以记录响应内容
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 执行下一个中间件
		c.Next()

		// 计算耗时
		latency := time.Since(startTime).Milliseconds()

		// 获取响应状态码和响应内容
		statusCode := c.Writer.Status()
		responseBody := blw.body.String()

		// 获取错误信息（如果有）
		var errMsg string
		if err := c.Errors.Last(); err != nil {
			errMsg = err.Error()
		}

		// 构建操作日志
		operationLog := &logger.OperationLog{
			UserID:      userID,
			Username:    username,
			IP:          clientIP,
			Method:      c.Request.Method,
			Path:        c.Request.URL.Path,
			QueryParams: queryParams,
			Body:        string(reqBody),
			Headers:     string(headersJSON),
			StatusCode:  statusCode,
			Response:    responseBody,
			Error:       errMsg,
			Latency:     latency,
			UserAgent:   c.Request.UserAgent(),
			Lang:        lang,
			Timezone:    timezone,
			OperateTime: time.Now(),
		}

		// 异步保存日志（避免阻塞请求）
		go func() {
			if err := logger.SaveOperationLog(operationLog); err != nil {
				logger.Error("Failed to save operation log:", err)
			}
		}()
	}
}

// GetClientIP 获取客户端真实IP
func GetClientIP(c *gin.Context) string {
	// 检查X-Forwarded-For头（通过代理的情况）
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return ip
	}
	// 检查X-Real-IP头
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	// 返回远程地址
	return c.ClientIP()
}

// ShouldLog 判断是否需要记录日志
func ShouldLog(path string, method string) bool {
	// 跳过健康检查和静态资源
	skipPaths := []string{
		"/health",
		"/favicon.ico",
	}

	for _, skipPath := range skipPaths {
		if skipPath == path {
			return false
		}
	}

	// 通常只记录非OPTIONS请求
	if method == http.MethodOptions {
		return false
	}

	return true
}