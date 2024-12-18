package model

import (
	"github.com/caarlos0/env"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log/slog"
	"os"
	"time"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type mLogger struct {
	*slog.Logger
}

func (l mLogger) Printf(message string, args ...any) {
	l.Info(message, args...)
}

var GormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // 表名使用单数, `User` -> `user`
	},
	DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束，必须手动创建或者在业务逻辑层维护
	Logger: logger.New(
		mLogger{Logger},
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	),
}

var NewDB *gorm.DB
var OriginalDB *gorm.DB

var Config struct {
	OriginalDbURL string `env:"Original_DB_URL"` // 原本的图床数据库地址
	NewDbURL      string `env:"New_DB_URL"`      // 要迁移的数据库地址
}

func Init() {
	err := env.Parse(&Config)
	if err != nil {
		panic(err)
	}
	originalDbSource := mysql.Open(Config.OriginalDbURL)
	newDbSource := mysql.Open(Config.NewDbURL)

	OriginalDB, err = gorm.Open(originalDbSource, GormConfig)
	if err != nil {
		panic(err)
	}
	err = OriginalDB.AutoMigrate(&OriginalImageTable{})
	if err != nil {
		panic(err)
	}

	NewDB, err = gorm.Open(newDbSource, GormConfig)
	if err != nil {
		panic(err)
	}
	err = NewDB.AutoMigrate(&NewImageTable{})
	if err != nil {
		panic(err)
	}

	slog.Info("database init success")

}
