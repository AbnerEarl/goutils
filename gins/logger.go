package gins

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogTmFmtWithMS = "2006-01-02 15:04:05.000"
)

func initLog(logDir, logfileName string) {
	l := &Log{
		LogDir:      logDir,
		LogFileName: logfileName,
		LogMinLevel: zap.InfoLevel,
		LogMaxSize:  1,
		MaxBackups:  3,
	}
	core := InitLogger(l)
	logger := zap.New(core)
	zap.ReplaceGlobals(logger)
}

type Log struct {
	LogDir      string
	LogFileName string
	Stdout      bool
	LogMaxSize  int
	LogMaxAge   int
	MaxBackups  int
	LogCompress bool
	LocalTime   bool
	EnableColor bool
	JsonFormat  bool
	LogMinLevel zapcore.Level
}

func InitLogger(l *Log) zapcore.Core {
	opts := []zapcore.WriteSyncer{
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(l.LogDir, l.LogFileName), // ⽇志⽂件路径
			MaxSize:    l.LogMaxSize,                           // 单位为MB,默认为512MB
			MaxAge:     l.LogMaxAge,                            // 文件最多保存多少天
			LocalTime:  l.LocalTime,                            // 采用本地时间
			Compress:   l.LogCompress,                          // 是否压缩日志
			MaxBackups: l.MaxBackups,                           // 最多保存多少分
		}),
	}

	if l.Stdout {
		opts = append(opts, zapcore.AddSync(os.Stdout))
	}

	syncWriter := zapcore.NewMultiWriteSyncer(opts...)

	// 自定义时间输出格式
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format(LogTmFmtWithMS) + "]")
	}
	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + level.CapitalString() + "]")
	}

	// 自定义文件：行号输出项
	customCallerEncoder := func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + caller.TrimmedPath() + "]")
	}

	encoderConf := zapcore.EncoderConfig{
		CallerKey:      "caller_line", // 打印文件名和行数
		LevelKey:       "level_name",
		MessageKey:     "msg",
		TimeKey:        "ts",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     customTimeEncoder,   // 自定义时间格式
		EncodeLevel:    customLevelEncoder,  // 小写编码器
		EncodeCaller:   customCallerEncoder, // 全路径编码器
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// level大写染色编码器
	if l.EnableColor {
		encoderConf.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// json 格式化处理
	if l.JsonFormat {
		return zapcore.NewCore(zapcore.NewJSONEncoder(encoderConf),
			syncWriter, zap.NewAtomicLevelAt(l.LogMinLevel))
	}

	return zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConf),
		syncWriter, zap.NewAtomicLevelAt(l.LogMinLevel))
}

func LogError(msg ...string) {
	for _, s := range msg {
		zap.L().Error(s)
	}
}

func LogInfo(msg ...string) {
	for _, s := range msg {
		zap.L().Info(s)
	}
}

func LogPanic(msg ...string) {
	for _, s := range msg {
		zap.L().Panic(s)
	}
}

func LogDebug(msg ...string) {
	for _, s := range msg {
		zap.L().Debug(s)
	}
}

func LogWarn(msg ...string) {
	for _, s := range msg {
		zap.L().Warn(s)
	}
}

func LogFatal(msg ...string) {
	for _, s := range msg {
		zap.L().Fatal(s)
	}
}
