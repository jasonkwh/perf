package perf

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"time"

	"go.uber.org/zap"
)

const pprofPort = 6060

type PProfServer struct {
	httpServer *http.Server
}

func NewPProfServer(zl *zap.Logger) (*PProfServer, error) {
	psrv := &PProfServer{}

	psrv.httpServer = &http.Server{
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: 5 * time.Second,   //nolint:mnd
		ReadTimeout:       30 * time.Second,  //nolint:mnd
		WriteTimeout:      30 * time.Second,  //nolint:mnd
		IdleTimeout:       120 * time.Second, //nolint:mnd
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", pprofPort))
	if err != nil {
		return nil, ErrFailedServerListenerInit(err)
	}

	go func() {
		if err := psrv.httpServer.Serve(lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zl.Error("pprof server failed", zap.Error(err))
		}
	}()

	zl.Info("pprof server started", zap.Int("port", pprofPort))

	return psrv, nil
}

func (psrv *PProfServer) Close() error {
	if psrv.httpServer != nil {
		return psrv.httpServer.Close()
	}

	return nil
}
