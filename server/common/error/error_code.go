package error

// ErrorCode 错误码枚举
type ErrorCode int

const (
	// 成功码
	SuccessCode ErrorCode = 0

	// 通用错误码 (10000-10099)
	UnknownError           ErrorCode = 10000
	InternalServerError    ErrorCode = 10001
	InvalidRequest         ErrorCode = 10002
	ValidationFailed       ErrorCode = 10003
	MethodNotAllowed       ErrorCode = 10004
	TooManyRequests        ErrorCode = 10005

	// 认证错误码 (10100-10199)
	Unauthorized           ErrorCode = 10100
	Forbidden              ErrorCode = 10101
	TokenExpired           ErrorCode = 10102
	TokenInvalid           ErrorCode = 10103
	TokenMissing           ErrorCode = 10104
	RefreshTokenInvalid    ErrorCode = 10105
	UserNotFound           ErrorCode = 10106
	PasswordIncorrect      ErrorCode = 10107
	AccountLocked          ErrorCode = 10108
	AccountDisabled        ErrorCode = 10109
	EmailNotVerified       ErrorCode = 10110

	// 业务错误码 (20000-29999)
	BusinessError          ErrorCode = 20000

	// 参数错误码 (30000-30099)
	InvalidParameter       ErrorCode = 30000
	MissingParameter       ErrorCode = 30001
	ParameterTypeError     ErrorCode = 30002

	// 数据库错误码 (40000-40099)
	DatabaseError          ErrorCode = 40000
	RecordNotFound         ErrorCode = 40001
	RecordExists           ErrorCode = 40002
	TransactionFailed      ErrorCode = 40003

	// Redis错误码 (50000-50099)
	RedisError             ErrorCode = 50000
)

// errorCodeMap 错误码到国际化key的映射
var errorCodeMap = map[ErrorCode]string{
	// 通用错误
	UnknownError:        "error.unknown",
	InternalServerError: "error.internal_server_error",
	InvalidRequest:      "error.invalid_request",
	ValidationFailed:    "error.validation_failed",
	MethodNotAllowed:    "error.method_not_allowed",
	TooManyRequests:     "error.too_many_requests",

	// 认证错误
	Unauthorized:        "error.unauthorized",
	Forbidden:           "error.forbidden",
	TokenExpired:        "auth.token_expired",
	TokenInvalid:        "auth.token_invalid",
	TokenMissing:        "auth.token_missing",
	RefreshTokenInvalid: "auth.refresh_token_invalid",
	UserNotFound:        "auth.user_not_found",
	PasswordIncorrect:   "auth.password_incorrect",
	AccountLocked:       "auth.account_locked",
	AccountDisabled:     "auth.account_disabled",
	EmailNotVerified:    "auth.email_not_verified",

	// 参数错误
	InvalidParameter:   "error.invalid_request",
	MissingParameter:   "error.invalid_request",
	ParameterTypeError: "error.invalid_request",

	// 数据库错误
	DatabaseError:     "error.internal_server_error",
	RecordNotFound:    "error.not_found",
	RecordExists:      "error.validation_failed",
	TransactionFailed: "error.internal_server_error",

	// Redis错误
	RedisError: "error.internal_server_error",
}

// GetMessageKey 获取错误码对应的国际化key
func (code ErrorCode) GetMessageKey() string {
	if key, ok := errorCodeMap[code]; ok {
		return key
	}
	return "error.unknown"
}

// IsSuccess 判断是否成功
func (code ErrorCode) IsSuccess() bool {
	return code == SuccessCode
}

// IsError 判断是否错误
func (code ErrorCode) IsError() bool {
	return code != SuccessCode
}