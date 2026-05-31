package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"SupCaller/common/database"

	"gorm.io/gorm"
)

// OperationLog 操作日志模型
type OperationLog struct {
	gorm.Model
	UserID      uint      `json:"user_id" gorm:"column:user_id;index"`
	Username    string    `json:"username" gorm:"column:username;size:100"`
	IP          string    `json:"ip" gorm:"column:ip;size:50"`
	Method      string    `json:"method" gorm:"column:method;size:10"`
	Path        string    `json:"path" gorm:"column:path;size:500"`
	QueryParams string    `json:"query_params" gorm:"column:query_params;type:text"`
	Body        string    `json:"body" gorm:"column:body;type:text"`
	Headers     string    `json:"headers" gorm:"column:headers;type:text"`
	StatusCode  int       `json:"status_code" gorm:"column:status_code"`
	Response    string    `json:"response" gorm:"column:response;type:text"`
	Error       string    `json:"error" gorm:"column:error;type:text"`
	Latency     int64     `json:"latency" gorm:"column:latency"`
	UserAgent   string    `json:"user_agent" gorm:"column:user_agent;size:500"`
	Lang        string    `json:"lang" gorm:"column:lang;size:20"`
	Timezone    string    `json:"timezone" gorm:"column:timezone;size:50"`
	OperateTime time.Time `json:"operate_time" gorm:"column:operate_time"`
}

// TableName 表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

// OperationLogConfig 操作日志配置
type OperationLogConfig struct {
	EnableDB     bool   // 是否启用数据库存储
	EnableFile   bool   // 是否启用文件存储
	FileDir      string // 文件存储目录
	FileMaxSize  int    // 单文件最大大小(MB)
	BackupCount  int    // 备份文件数量
}

var operationLogConfig = OperationLogConfig{
	EnableDB:    true,
	EnableFile:  true,
	FileDir:     "./logs/operation",
	FileMaxSize: 10,
	BackupCount: 7,
}

// SetOperationLogConfig 设置操作日志配置
func SetOperationLogConfig(config OperationLogConfig) {
	operationLogConfig = config
}

// InitOperationLog 初始化操作日志（创建表）
func InitOperationLog() error {
	if operationLogConfig.EnableDB {
		err := database.DB.AutoMigrate(&OperationLog{})
		if err != nil {
			return fmt.Errorf("failed to migrate operation_logs table: %v", err)
		}
	}

	if operationLogConfig.EnableFile {
		err := os.MkdirAll(operationLogConfig.FileDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}
	}

	return nil
}

// SaveOperationLog 保存操作日志
func SaveOperationLog(log *OperationLog) error {
	log.OperateTime = time.Now()

	var err error

	// 保存到数据库
	if operationLogConfig.EnableDB {
		err = database.DB.Create(log).Error
		if err != nil {
			Error("Failed to save operation log to database:", err)
		}
	}

	// 保存到文件
	if operationLogConfig.EnableFile {
		fileErr := saveOperationLogToFile(log)
		if fileErr != nil {
			Error("Failed to save operation log to file:", fileErr)
			if err == nil {
				err = fileErr
			}
		}
	}

	return err
}

// saveOperationLogToFile 保存操作日志到文件
func saveOperationLogToFile(log *OperationLog) error {
	filePath := getOperationLogFilePath()
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 检查文件大小，需要轮转
	if err := rotateOperationLogFile(filePath); err != nil {
		return err
	}

	// 格式化日志内容
	logJSON, err := json.Marshal(log)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(logJSON) + "\n")
	return err
}

// getOperationLogFilePath 获取操作日志文件路径
func getOperationLogFilePath() string {
	date := time.Now().Format("2006-01-02")
	return filepath.Join(operationLogConfig.FileDir, fmt.Sprintf("operation_%s.log", date))
}

// rotateOperationLogFile 轮转日志文件
func rotateOperationLogFile(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// 检查文件大小是否超过限制（MB）
	maxSizeBytes := int64(operationLogConfig.FileMaxSize) * 1024 * 1024
	if fileInfo.Size() < maxSizeBytes {
		return nil
	}

	// 轮转文件
	for i := operationLogConfig.BackupCount - 1; i >= 0; i-- {
		srcPath := filePath
		if i > 0 {
			srcPath = fmt.Sprintf("%s.%d", filePath, i)
		}
		dstPath := fmt.Sprintf("%s.%d", filePath, i+1)

		if _, err := os.Stat(srcPath); err == nil {
			if err := os.Rename(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	// 创建新文件
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	file.Close()

	return nil
}

// QueryOperationLogs 查询操作日志
func QueryOperationLogs(query OperationLogQuery) ([]OperationLog, int64, error) {
	var logs []OperationLog
	var count int64

	db := database.DB.Model(&OperationLog{})

	if query.UserID != 0 {
		db = db.Where("user_id = ?", query.UserID)
	}
	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.IP != "" {
		db = db.Where("ip = ?", query.IP)
	}
	if query.Method != "" {
		db = db.Where("method = ?", query.Method)
	}
	if query.Path != "" {
		db = db.Where("path LIKE ?", "%"+query.Path+"%")
	}
	if query.StartTime != nil {
		db = db.Where("operate_time >= ?", query.StartTime)
	}
	if query.EndTime != nil {
		db = db.Where("operate_time <= ?", query.EndTime)
	}

	// 统计总数
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	err = db.Order("operate_time DESC").
		Offset(query.Offset).
		Limit(query.Limit).
		Find(&logs).Error

	return logs, count, err
}

// OperationLogQuery 查询条件
type OperationLogQuery struct {
	UserID    uint
	Username  string
	IP        string
	Method    string
	Path      string
	StartTime *time.Time
	EndTime   *time.Time
	Offset    int
	Limit     int
}

// GetOperationLogByID 根据ID获取操作日志
func GetOperationLogByID(id uint) (*OperationLog, error) {
	var log OperationLog
	err := database.DB.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// DeleteOperationLog 删除操作日志
func DeleteOperationLog(id uint) error {
	return database.DB.Delete(&OperationLog{}, id).Error
}

// CleanOldLogs 清理旧日志
func CleanOldLogs(days int) error {
	if !operationLogConfig.EnableDB {
		return nil
	}

	expireTime := time.Now().AddDate(0, 0, -days)
	return database.DB.Where("operate_time < ?", expireTime).Delete(&OperationLog{}).Error
}