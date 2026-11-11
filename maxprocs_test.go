package perf_test

import (
	"runtime"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/jasonkwh/perf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestAutoMaxProcs_HonorsEnv(t *testing.T) {
	as := assert.New(t)
	t.Setenv("GOMAXPROCS", "2")

	core, logs := observer.New(zapcore.DebugLevel)
	zl := zap.New(core)

	prev := runtime.GOMAXPROCS(0)
	t.Cleanup(func() { runtime.GOMAXPROCS(prev) })

	perf.AutoMaxProcs(zl)

	entries := logs.AllUntimed()

	if as.Len(entries, 1) {
		as.Equal(zapcore.DebugLevel, entries[0].Level)
		as.Equal(`maxprocs: Honoring GOMAXPROCS="2" as set in environment`, entries[0].Message)
	}
}
