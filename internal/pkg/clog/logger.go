package clog

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

type Clog struct {
	prefix string
	*slog.Logger
}

const chatCtxKey = "chat"

type chatCtxMap map[string]interface{}

func getChatCtx(ctx context.Context) chatCtxMap {
	cmp := ctx.Value(chatCtxKey)
	if cmp == nil {
		return nil
	}
	return cmp.(chatCtxMap)
}

func getTraceID(ctx context.Context) string {
	c := getChatCtx(ctx)
	if c == nil {
		return ""
	}
	return c["traceId"].(string)
}

func getUserID(ctx context.Context) string {
	c := getChatCtx(ctx)
	if c == nil {
		return ""
	}
	return c["userId"].(string)
}

var ModuleLogger sync.Map

func NewClog(prefix string, options *Options) *Clog {
	logger := &Clog{
		prefix: prefix,
	}
	configLogger(logger, options)
	return logger
}

func configLogger(clog *Clog, options *Options) {
	opts := &slog.HandlerOptions{
		Level:     options.Level,
		AddSource: true,
	}
	var err error
	b := os.Stdout
	if options.Dir != "" {
		b, err = os.OpenFile(options.Dir, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic("Failed to open")
		}
	}

	var handler slog.Handler = slog.NewJSONHandler(b, opts)
	clog.Logger = slog.New(handler)
}

func RegisterLogger(moduleName string, options *Options) *Clog {
	var clog *Clog
	l, ok := ModuleLogger.Load(moduleName)
	if !ok {
		clog = NewClog(moduleName, options)
		ModuleLogger.Store(moduleName, clog)
	} else {
		clog = l.(*Clog)
	}
	return clog
}

func (c *Clog) InfoContext(ctx context.Context, msg string, args ...any) {
	tid := getTraceID(ctx)
	uid := getUserID(ctx)
	args = append(args, slog.String("traceId", tid), slog.String("userId", uid), slog.String("module", c.prefix))
	c.Logger.InfoContext(ctx, msg, args...)
}

func (c *Clog) DebugContext(ctx context.Context, msg string, args ...any) {
	tid := getTraceID(ctx)
	uid := getUserID(ctx)
	args = append(args, slog.String("traceId", tid), slog.String("userId", uid), slog.String("module", c.prefix))
	c.Logger.DebugContext(ctx, msg, args...)
}

func (c *Clog) WarnContext(ctx context.Context, msg string, args ...any) {
	tid := getTraceID(ctx)
	uid := getUserID(ctx)
	args = append(args, slog.String("traceId", tid), slog.String("userId", uid), slog.String("module", c.prefix))
	c.Logger.WarnContext(ctx, msg, args...)
}

func (c *Clog) ErrorContext(ctx context.Context, msg string, args ...any) {
	tid := getTraceID(ctx)
	uid := getUserID(ctx)
	args = append(args, slog.String("traceId", tid), slog.String("userId", uid), slog.String("module", c.prefix))
	c.Logger.ErrorContext(ctx, msg, args...)
}
