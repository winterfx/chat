package clog

import (
	"context"
	"log/slog"
	"testing"
)

func TestInfo(t *testing.T) {
	logger := RegisterLogger("api-server", NewOptions(WithLevel(slog.LevelInfo)))
	logger.Info("this is info")
	logger.InfoContext(context.Background(), "this is info")
}
