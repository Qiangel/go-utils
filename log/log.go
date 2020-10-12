package log

import (
	"os"
	"time"

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
)

var logger *log

const (
	FileName   = "./logs/"
	MaxSize    = 10
	MaxBackups = 30
	MaxAge     = 7
)
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
	PanicLevel = "panic"
)

type log struct {
	opts options
}

var defaultOptions = options{
	level:          DebugLevel,
	enableSaveFile: false,
	enableSaveDB:   false,
	caller:         true,
	fileName:       FileName + defaultFileName(),
	maxSize:        MaxSize,
	maxBackups:     MaxBackups,
	maxAge:         MaxAge,
	compress:       true,
}

type options struct {
	version        string //版本
	level          string //日志级别
	enableSaveFile bool   //是否开启文件存储
	enableSaveDB   bool   //是否开启存储到数据库
	caller         bool   //显示行号等
	fileName       string //存储路径和文件
	maxSize        int    //每个文件的最大尺寸 M
	maxBackups     int    //最多保存多少个备份
	maxAge         int    //文件最多保存多少天
	compress       bool   //是否压缩

	//logger
	logger *zap.Logger
	//sugar logger
	sugaredLogger *zap.SugaredLogger
}

type Option func(*options)

//WithVersion 版本
func WithVersion(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

//WithEnableSaveToFile 是否开启存储到文件
func WithEnableSaveFile(enable bool) Option {
	return func(o *options) {
		o.enableSaveFile = enable
	}
}

//WithEnableSaveDB 是否开启存储到数据库
func WithEnableSaveDB(enable bool) Option {
	return func(o *options) {
		o.enableSaveDB = enable
	}
}

//WithFilePath 路径
func WithFilePath(filePath string) Option {
	return func(o *options) {
		o.fileName = filePath + defaultFileName()
	}
}

//WithMaxSize 文件大小单位m
func WithMaxSize(maxSize int) Option {
	return func(o *options) {
		o.maxSize = maxSize
	}
}

//WithMaxBackups 文件个数
func WithMaxBackups(maxBackups int) Option {
	return func(o *options) {
		o.maxBackups = maxBackups
	}
}

//WithMaxAge 保存最大天数
func WithMaxAge(maxAge int) Option {
	return func(o *options) {
		o.maxAge = maxAge
	}
}

//WithCompress 是否开启压缩
func WithCompress(compress bool) Option {
	return func(o *options) {
		o.compress = compress
	}
}

func Init(opts ...Option) {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	build(&o)
	logger = &log{opts: o}
}

func build(opts *options) {
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   opts.fileName,
		MaxSize:    opts.maxSize,
		MaxBackups: opts.maxBackups,
		MaxAge:     opts.maxAge,
		Compress:   opts.compress,
	})
	var level zapcore.Level
	switch opts.level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "panic":
		level = zap.PanicLevel
	default:
		level = zap.InfoLevel
	}
	conf := zap.NewProductionEncoderConfig()
	conf.TimeKey = "time"
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.LineEnding = zapcore.DefaultLineEnding
	conf.EncodeDuration = zapcore.SecondsDurationEncoder
	conf.EncodeCaller = zapcore.FullCallerEncoder //全路径编码器 绝对路径？
	enc := zapcore.NewJSONEncoder(conf)

	//日志输出方式
	var wss []zapcore.WriteSyncer
	//输出到控制台
	wss = append(wss, zapcore.AddSync(os.Stdout))
	if opts.enableSaveFile {
		wss = append(wss, ws) //日志切割文件
	}
	if opts.enableSaveDB {
		wss = append(wss, zapcore.AddSync(hk)) //日志存入数据库
	}
	writeSs := zapcore.NewMultiWriteSyncer(wss...)
	core := zapcore.NewCore(
		enc,
		writeSs,
		level,
	)
	//zap option 参数
	var zapOpts []zap.Option
	if opts.version != "" {
		zapOpts = append(zapOpts, zap.Fields(zap.String("version", opts.version)))
	}
	//todo trace id

	opts.logger = zap.New(
		core,
		zapOpts...,
	)
	if opts.caller {
		opts.logger = opts.logger.WithOptions(
			zap.Development(),    //开启文件和行号
			zap.AddCaller(),      //开启开发模式，堆栈跟踪
			zap.AddCallerSkip(1), //跳
		)
	}
	//sugaredLogger 构建
	opts.sugaredLogger = opts.logger.Sugar()
}

//flush buffered
func Sync() {
	if logger == nil {
		return
	}
	_ = logger.opts.sugaredLogger.Sync()
}

//defaultFileName 生成文件名称
func defaultFileName() string {
	fileName := time.Now().Format("20060102") + ".log"
	return fileName
}
