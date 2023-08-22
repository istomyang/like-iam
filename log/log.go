package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

// HookBody is currently used to server as hook body give the hook function.
type HookBody = zapcore.Entry

// HookFunc defines a hook function developers register to log.
type HookFunc = func(c *context.Context, body HookBody) error

type Logger struct {
	inner *zap.Logger
	ctx   context.Context
	level Level
	// I temporarily use zap to manage hooks feature.
	// hooks map[Level][]HookFunc
}

func NewLogger(ctx context.Context, options *Options) (r *Logger) {
	l, err := options.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to build logger: %v", err))
	}
	r = &Logger{inner: l, ctx: ctx, level: DebugLevel.into(options.Level)}
	return
}

const (
	// DebugEnabledKey is a Context key to enable/disable debug-level logging in this context.
	DebugEnabledKey = "DebugEnabledKey"
	// XRequestIDKey is a convenient key to filter log infos with an independent caller-procedure.
	XRequestIDKey = "X-Request-ID"
	// UserNameKey is a name of user.
	UserNameKey = "UserName"
)

// L will extract values from context, adding to common logger fields and returning a clone logger.
func (t *Logger) L(ctx context.Context) *Logger {
	var inner = t.inner

	if requestId := ctx.Value(XRequestIDKey); requestId != nil {
		inner = inner.With(zap.Any(XRequestIDKey, requestId))
	}
	if username := ctx.Value(UserNameKey); username != nil {
		inner = inner.With(zap.Any(UserNameKey, username))
	}

	var level Level
	if ctx.Value(DebugEnabledKey).(bool) {
		level = DebugLevel
	}

	return &Logger{
		inner: inner,
		ctx:   ctx,
		level: level,
	}
}

func (t *Logger) Debug(msg string, fields ...Field) {
	if !t.allow(DebugLevel) {
		return
	}
	var zapFields = make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.toZapField()
	}
	t.inner.Debug(msg, zapFields...)
}
func (t *Logger) Debugf(format string, v ...any) {
	if !t.allow(DebugLevel) {
		return
	}
	t.inner.Sugar().Debugf(format, v...)
}
func (t *Logger) Debugw(msg string, kv ...any) {
	if !t.allow(DebugLevel) {
		return
	}
	t.inner.Sugar().Debugw(msg, kv...)
}
func (t *Logger) Info(msg string, fields ...Field) {
	if !t.allow(InfoLevel) {
		return
	}
	var zapFields = make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.toZapField()
	}
	t.inner.Info(msg, zapFields...)
}
func (t *Logger) Infof(format string, v ...any) {
	if !t.allow(InfoLevel) {
		return
	}
	t.inner.Sugar().Infof(format, v...)
}
func (t *Logger) Infow(msg string, kv ...any) {
	if !t.allow(InfoLevel) {
		return
	}
	t.inner.Sugar().Infow(msg, kv...)
}
func (t *Logger) Warn(msg string, fields ...Field) {
	if !t.allow(WarnLevel) {
		return
	}
	var zapFields = make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.toZapField()
	}
	t.inner.Warn(msg, zapFields...)
}
func (t *Logger) Warnf(format string, v ...any) {
	if !t.allow(WarnLevel) {
		return
	}
	t.inner.Sugar().Warnf(format, v...)
}
func (t *Logger) Warnw(msg string, kv ...any) {
	if !t.allow(WarnLevel) {
		return
	}
	t.inner.Sugar().Warnw(msg, kv...)
}
func (t *Logger) Error(msg string, fields ...Field) {
	if !t.allow(ErrorLevel) {
		return
	}
	var zapFields = make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.toZapField()
	}
	t.inner.Error(msg, zapFields...)
}
func (t *Logger) Errorf(format string, v ...any) {
	if !t.allow(ErrorLevel) {
		return
	}
	t.inner.Sugar().Errorf(format, v...)
}
func (t *Logger) Errorw(msg string, kv ...any) {
	if !t.allow(ErrorLevel) {
		return
	}
	t.inner.Sugar().Errorw(msg, kv...)
}
func (t *Logger) Panic(msg string, fields ...Field) {
	if !t.allow(PanicLevel) {
		return
	}
	var zapFields = make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.toZapField()
	}
	t.inner.Panic(msg, zapFields...)
}
func (t *Logger) Panicf(format string, v ...any) {
	if !t.allow(PanicLevel) {
		return
	}
	t.inner.Sugar().Panicf(format, v...)
}
func (t *Logger) Panicw(msg string, kv ...any) {
	if !t.allow(PanicLevel) {
		return
	}
	t.inner.Sugar().Panicw(msg, kv...)
}
func (t *Logger) Fatal(msg string, fields ...Field) {
	if !t.allow(FatalLevel) {
		return
	}
	var zapFields = make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = f.toZapField()
	}
	t.inner.Fatal(msg, zapFields...)
}
func (t *Logger) Fatalf(format string, v ...any) {
	if !t.allow(FatalLevel) {
		return
	}
	t.inner.Sugar().Fatalf(format, v...)
}
func (t *Logger) Fatalw(msg string, kv ...any) {
	if !t.allow(FatalLevel) {
		return
	}
	t.inner.Sugar().Fatalw(msg, kv...)
}

// V is used to re-define those log-function's level.
func (t *Logger) V(l Level) *LevelLogger {
	return &LevelLogger{inner: t.inner, level: l}
}

func (t *Logger) Sync() {
	_ = t.inner.Sync()
}

func (t *Logger) allow(l Level) bool {
	// I put context "enable-debug-key" into t.level.
	return t.level.allow(l)
}

// RegisterHooks should register once, multiple times will replace previous registered hooks
func (t *Logger) RegisterHooks(l Level, fff ...HookFunc) {
	//h := append(t.hooks[l], f)
	//n := make([]HookFunc, len(h))
	//for i, hookFunc := range h {
	//	n[i] = hookFunc
	//}
	//t.hooks[l] = n

	ff := make([]func(zapcore.Entry) error, len(fff))
	for i, f := range fff {
		ff[i] = func(entry zapcore.Entry) error {
			return f(&t.ctx, entry)
		}
	}
	t.inner = t.inner.WithOptions(zap.Hooks(ff...))
}

type LevelLogger struct {
	inner *zap.Logger
	level Level
}

func (t *LevelLogger) Info(msg string, fields ...Field) {
	if entry := t.inner.Check(t.level.intoZapLevel(), msg); entry != nil {
		var zapFields = make([]zap.Field, len(fields))
		for i, f := range fields {
			zapFields[i] = f.toZapField()
		}
		entry.Write(zapFields...)
	}
}
func (t *LevelLogger) Infof(format string, v ...any) {
	if entry := t.inner.Check(t.level.intoZapLevel(), fmt.Sprint(format, v)); entry != nil {
		entry.Write()
	}
}
func (t *LevelLogger) Infow(msg string, kv ...any) {
	if entry := t.inner.Check(t.level.intoZapLevel(), msg); entry != nil {
		if len(kv) == 0 {
			entry.Write()
		}
		if len(kv)%2 != 0 {
			panic(fmt.Sprintf("args is not an even number of arguments: %v", kv))
		}

		fields := make([]zap.Field, len(kv)/2)

		for i := 0; i < len(kv); i += 2 {
			key, val := kv[i], kv[i+1]
			if ketStr, ok := key.(string); ok {
				fields[i] = zap.Any(ketStr, val)
			} else {
				panic(fmt.Sprintf("key is not a string: %v", key))
			}
		}

		entry.Write(fields...)
	}
}

type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

func (l Level) string() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

func (l Level) into(s string) Level {
	switch s {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	case "panic":
		return PanicLevel
	default:
		return PanicLevel
	}
}

func (l Level) intoZapLevel() zapcore.Level {
	r, _ := zapcore.ParseLevel(l.string())
	return r
}

func (l Level) allow(t Level) bool {
	return t >= l
}

type Field struct {
	key   string
	value any
}

func (f Field) toZapField() zap.Field {
	return zap.Any(f.key, f.value)
}

func Any(key string, value any) Field {
	return Field{key: key, value: value}
}

var (
	std = NewLogger(context.Background(), NewOptions("global", nil))
	mut sync.Mutex
)

func Default() *Logger { return std }

func Init(ctx context.Context, options *Options) {
	mut.Lock()
	defer mut.Unlock()
	std = NewLogger(ctx, options)
}

func Debug(msg string, fields ...Field) {
	std.Debug(msg, fields...)
}
func Debugf(format string, v ...any) {
	std.Debugf(format, v...)
}
func Debugw(msg string, kv ...any) {
	std.Debugw(msg, kv)
}
func Info(msg string, fields ...Field) {
	std.Info(msg, fields...)
}
func Infof(format string, v ...any) {
	std.Infof(format, v)
}
func Infow(msg string, kv ...any) {
	std.Infow(msg, kv...)
}
func Warn(msg string, fields ...Field) {
	std.Warn(msg, fields...)
}
func Warnf(format string, v ...any) {
	std.Warnf(format, v...)
}
func Warnw(msg string, kv ...any) {
	std.Warnw(msg, kv...)
}
func Error(msg string, fields ...Field) {
	std.Error(msg, fields...)
}
func Errorf(format string, v ...any) {
	std.Errorf(format, v)
}
func Errorw(msg string, kv ...any) {
	std.Errorw(msg, kv...)
}
func Panic(msg string, fields ...Field) {
	std.Panic(msg, fields...)
}
func Panicf(format string, v ...any) {
	std.Panicf(format, v)
}
func Panicw(msg string, kv ...any) {
	std.Panicw(msg, kv...)
}
func Fatal(msg string, fields ...Field) {
	std.Fatal(msg, fields...)
}
func Fatalf(format string, v ...any) {
	std.Fatalf(format, v...)
}
func Fatalw(msg string, kv ...any) {
	std.Fatalw(msg, kv...)
}

func Sync() {
	std.Sync()
}

func L(ctx context.Context) *Logger {
	return std.L(ctx)
}

func V(l Level) *LevelLogger {
	return std.V(l)
}

func RegisterHooks(l Level, fff ...HookFunc) {
	std.RegisterHooks(l, fff...)
}
