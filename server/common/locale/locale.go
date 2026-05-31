package locale

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

const (
	// 默认语言
	DefaultLang = "zh-CN"
	// 默认时区
	DefaultTimezone = "Asia/Shanghai"
	// 语言请求头
	LangHeader = "Accept-Language"
	// 时区请求头
	TimezoneHeader = "X-Timezone"
)

var (
	translations   = make(map[string]map[string]string)
	timezoneCache  = make(map[string]*time.Location)
	mu             sync.RWMutex
	loaded         bool
)

// InitLocale 初始化国际化
func InitLocale() error {
	mu.Lock()
	defer mu.Unlock()

	if loaded {
		return nil
	}

	// 加载语言文件
	localeDir := "./common/locale/data"
	files, err := os.ReadDir(localeDir)
	if err != nil {
		return fmt.Errorf("failed to read locale directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := filepath.Ext(file.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		lang := file.Name()[:len(file.Name())-len(ext)]
		filePath := filepath.Join(localeDir, file.Name())

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read locale file %s: %v", filePath, err)
		}

		var data map[string]string
		if err := yaml.Unmarshal(content, &data); err != nil {
			return fmt.Errorf("failed to parse locale file %s: %v", filePath, err)
		}

		translations[lang] = data
	}

	loaded = true
	return nil
}

// GetLangFromContext 从请求上下文获取语言
func GetLangFromContext(c *gin.Context) string {
	lang := c.GetHeader(LangHeader)
	if lang == "" {
		lang = DefaultLang
	}
	return normalizeLang(lang)
}

// GetTimezoneFromContext 从请求上下文获取时区
func GetTimezoneFromContext(c *gin.Context) *time.Location {
	tzStr := c.GetHeader(TimezoneHeader)
	if tzStr == "" {
		tzStr = DefaultTimezone
	}

	return GetTimezone(tzStr)
}

// GetTimezone 获取时区
func GetTimezone(tzStr string) *time.Location {
	mu.RLock()
	if loc, ok := timezoneCache[tzStr]; ok {
		mu.RUnlock()
		return loc
	}
	mu.RUnlock()

	loc, err := time.LoadLocation(tzStr)
	if err != nil {
		loc, _ = time.LoadLocation(DefaultTimezone)
	}

	mu.Lock()
	timezoneCache[tzStr] = loc
	mu.Unlock()

	return loc
}

// Translate 翻译文本
func Translate(lang, key string) string {
	return TranslateWithArgs(lang, key)
}

// TranslateWithArgs 翻译文本（带参数）
func TranslateWithArgs(lang, key string, args ...interface{}) string {
	mu.RLock()
	defer mu.RUnlock()

	lang = normalizeLang(lang)

	// 查找语言包
	if langMap, ok := translations[lang]; ok {
		if value, ok := langMap[key]; ok {
			if len(args) > 0 {
				return fmt.Sprintf(value, args...)
			}
			return value
		}
	}

	// 如果当前语言找不到，尝试默认语言
	if lang != DefaultLang {
		return TranslateWithArgs(DefaultLang, key, args...)
	}

	// 如果都找不到，返回key本身
	return key
}

// TranslateFromContext 从上下文翻译
func TranslateFromContext(c *gin.Context, key string) string {
	lang := GetLangFromContext(c)
	return Translate(lang, key)
}

// TranslateFromContextWithArgs 从上下文翻译（带参数）
func TranslateFromContextWithArgs(c *gin.Context, key string, args ...interface{}) string {
	lang := GetLangFromContext(c)
	return TranslateWithArgs(lang, key, args...)
}

// FormatTime 根据时区格式化时间
func FormatTime(c *gin.Context, t time.Time, layout string) string {
	loc := GetTimezoneFromContext(c)
	return t.In(loc).Format(layout)
}

// FormatTimeWithTimezone 根据指定时区格式化时间
func FormatTimeWithTimezone(t time.Time, tz *time.Location, layout string) string {
	return t.In(tz).Format(layout)
}

// Now 获取当前时间（带时区）
func Now(c *gin.Context) time.Time {
	loc := GetTimezoneFromContext(c)
	return time.Now().In(loc)
}

// normalizeLang 标准化语言代码
func normalizeLang(lang string) string {
	switch lang {
	case "zh", "zh-CN", "zh-cn", "zh_Hans", "zh-hans":
		return "zh-cn"
	case "en", "en-US", "en-us", "en_GB", "en-gb":
		return "en-us"
	default:
		return DefaultLang
	}
}

// GetSupportedLangs 获取支持的语言列表
func GetSupportedLangs() []string {
	mu.RLock()
	defer mu.RUnlock()

	langs := make([]string, 0, len(translations))
	for lang := range translations {
		langs = append(langs, lang)
	}
	return langs
}

// HasLang 检查是否支持指定语言
func HasLang(lang string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := translations[normalizeLang(lang)]
	return ok
}