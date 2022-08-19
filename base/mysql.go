package base

import (
	"context"
	"fmt"
	"github.com/derekAHua/goLib/env"
	"github.com/derekAHua/goLib/utils"
	"github.com/derekAHua/goLib/zlog"
	"go.uber.org/zap"

	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	ormUtil "gorm.io/gorm/utils"
)

type MysqlConf struct {
	Service         string        `yaml:"service"`
	DataBase        string        `yaml:"database"`
	Addr            string        `yaml:"addr"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Charset         string        `yaml:"charset"`
	MaxIdleConns    int           `yaml:"maxIdleConns"`
	MaxOpenConns    int           `yaml:"maxOpenConns"`
	ConnMaxLifeTime time.Duration `yaml:"connMaxLifeTime"`
	ConnTimeOut     time.Duration `yaml:"connTimeOut"`
	WriteTimeOut    time.Duration `yaml:"writeTimeOut"`
	ReadTimeOut     time.Duration `yaml:"readTimeOut"`
}

func (conf *MysqlConf) checkConf() {
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = 10
	}
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = 1000
	}
	if conf.ConnMaxLifeTime == 0 {
		conf.ConnMaxLifeTime = 3600 * time.Second
	}
	if conf.ConnTimeOut == 0 {
		conf.ConnTimeOut = 3 * time.Second
	}
	if conf.WriteTimeOut == 0 {
		conf.WriteTimeOut = 1 * time.Second
	}
	if conf.ReadTimeOut == 0 {
		conf.ReadTimeOut = 1 * time.Second
	}
}

func InitMysqlClient(conf MysqlConf) (client *gorm.DB, err error) {
	conf.checkConf()

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=True&loc=Asia%%2FShanghai",
		conf.User,
		conf.Password,
		conf.Addr,
		conf.DataBase,
		conf.ConnTimeOut,
		conf.ReadTimeOut,
		conf.WriteTimeOut)

	if conf.Charset != "" {
		dsn = dsn + "&charset=" + conf.Charset
	}

	c := &gorm.Config{
		SkipDefaultTransaction:                   true,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   newLogger(&conf),
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		AllowGlobalUpdate:                        false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	}

	client, err = gorm.Open(mysql.Open(dsn), c)
	if err != nil {
		return client, err
	}

	sqlDB, err := client.DB()
	if err != nil {
		return client, err
	}

	// 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)

	// 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

	// 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(conf.ConnMaxLifeTime)

	return client, nil
}

func newLogger(conf *MysqlConf) logger.Interface {
	s := conf.Service
	if conf.Service == "" {
		s = conf.DataBase
	}

	return &ormLogger{
		Service:  s,
		Addr:     conf.Addr,
		Database: conf.DataBase,
	}
}

type ormLogger struct {
	Service  string
	Addr     string
	Database string
}

// LogMode log mode
func (l *ormLogger) LogMode(_ logger.LogLevel) logger.Interface {
	newLogger := *l
	return &newLogger
}

// Info print info.
func (l ormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormUtil.FileWithLineNum()}, data...)...)
	zlog.InfoLogger(nil, zlog.LogNameMysql, m, l.commonFields(ctx)...)
}

// Warn print warn messages.
func (l ormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormUtil.FileWithLineNum()}, data...)...)
	zlog.WarnLogger(nil, zlog.LogNameMysql, m, l.commonFields(ctx)...)
}

// Error print error messages.
func (l ormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	m := fmt.Sprintf(msg, append([]interface{}{ormUtil.FileWithLineNum()}, data...)...)
	zlog.ErrorLogger(nil, zlog.LogNameMysql, m, l.commonFields(ctx)...)
}

func (l ormLogger) commonFields(ctx context.Context) []zlog.Field {
	var logID, requestID string
	if c, ok := ctx.(*gin.Context); ok && c != nil {
		logID, _ = ctx.Value(zlog.ContextKeyLogId).(string)
		requestID, _ = ctx.Value(zlog.ContextKeyRequestId).(string)
	}

	fields := []zlog.Field{
		zlog.WithTopicField(zlog.LogNameMysql),
		zap.String("logId", logID),
		zap.String("requestId", requestID),
		zap.String("protobuf", "mysql"),
		zap.String("module", env.GetAppName()),
		zap.String("service", l.Service),
		zap.String("addr", l.Addr),
		zap.String("db", l.Database),
	}
	return fields
}

// Trace print sql message.
func (l ormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	end := time.Now()
	elapsed := end.Sub(begin)
	cost := float64(elapsed.Nanoseconds()/1e4) / 100.0

	// 请求是否成功
	msg := "mysql do success"
	ralCode := -0
	if err != nil {
		msg = err.Error()
		ralCode = -1
	}

	sql, rows := fc()
	fileLineNum := ormUtil.FileWithLineNum()

	fields := l.commonFields(ctx)
	fields = append(fields,
		zap.String("sql", sql),
		zap.Int64("affectedRow", rows),
		zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
		zap.String("requestStartTime", utils.GetFormatRequestTime(begin)),
		zap.String("fileLine", fileLineNum),
		zap.Float64("cost", cost),
		zap.Int("ralCode", ralCode),
	)

	zlog.InfoLogger(nil, zlog.LogNameMysql, msg, fields...)
}
