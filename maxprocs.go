package perf

import (
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

func AutoMaxProcs(zl *zap.Logger) {
	if _, err := maxprocs.Set(maxprocs.Logger(zl.Sugar().Debugf)); err != nil {
		zl.Debug("automaxprocs initialization failed", zap.Error(err))
	}
}
