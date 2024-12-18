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

var DB *gorm.DB
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

var Config struct {
	DbURL    string `env:"DB_URL"`
	HostName string `env:"HOST_NAME" envDefault:"localhost:8000"`
}

func Init() {
	err := env.Parse(&Config)
	if err != nil {
		panic(err)
	}
	source := mysql.Open(Config.DbURL)
	DB, err = gorm.Open(source, GormConfig)
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&ImageTable{})
	if err != nil {
		panic(err)
	}

}
