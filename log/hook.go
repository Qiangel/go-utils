package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Qiangel/go-utils/mysql"

	"gorm.io/gorm"
)

var hk *Hooker

type Hooker struct {
	db *gorm.DB
}

func (h *Hooker) Write(p []byte) (n int, err error) {
	var hookLog HookLog
	err = json.Unmarshal(p, &hookLog)
	if err != nil {
		return
	}
	fmt.Println("hookLog:", hookLog)
	//写入数据库
	item := &LogItem{
		Level:      hookLog.Level,
		Message:    hookLog.Msg,
		TraceID:    hookLog.TraceID,
		UserID:     hookLog.UserID,
		Caller:     hookLog.Caller,
		Data:       string(p),
		ServerName: hookLog.ServerName,
		Version:    hookLog.Version,
	}
	if err = h.db.Create(&item).Error; err != nil {
		return 0, err
	}
	return
}

type HookLog struct {
	Level      string `json:"level"`
	Time       string `json:"time"`
	Caller     string `json:"caller"`
	Msg        string `json:"msg"`
	Version    string `json:"version"`
	ServerName string `json:"server_name"`
	TraceID    string `json:"trace_id"`
	UserID     string `json:"user_id"`
}

func InitHook(config mysql.Config) (err error) {
	db, err := mysql.Open(config)
	if err != nil {
		return err
	}
	db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	//开始写入数据库
	if db.Migrator().CurrentDatabase() != config.Dbname {
		return errors.New("some db err") //不可能进入
	}

	err = db.AutoMigrate(new(LogItem))
	if err != nil {
		return
	}
	hk = &Hooker{db: db}
	return
}

type LogItem struct {
	ID         uint      `gorm:"column:id;primary_key;auto_increment;"`                              // id
	ServerName string    `gorm:"column:server_name;size:50;not null;default:'';comment:服务名称;index;"` // 服务名称
	Level      string    `gorm:"column:level;size:20;not null;default:'';comment:日志级别;index;"`       // 日志级别
	Message    string    `gorm:"column:message;size:1024;not null;default:'';comment:消息;"`           // 消息
	TraceID    string    `gorm:"column:trace_id;size:128;not null;default:'';comment:跟踪id;index;"`   // 跟踪ID
	UserID     string    `gorm:"column:user_id;size:36;not null;default:'';comment:用户id;index;"`     // 用户ID
	Caller     string    `gorm:"column:caller;size:256;not null;default:'';comment:caller信息;"`       // Caller
	Data       string    `gorm:"column:data;type:text;not null;comment:日志数据;"`                       // 日志数据(json)
	Version    string    `gorm:"column:version;index;size:32;not null;default:'';comment:版本号;"`      // 服务版本号
	CreatedAt  time.Time `gorm:"column:created_at;index"`                                            // 创建时间
}

// TableName 表名
func (LogItem) TableName() string {
	return "log_item"
}
