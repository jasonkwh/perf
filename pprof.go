package perf

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"
)

const pprofPort = 6060

type PProfServer struct {
	httpServer *http.Server
	lis        net.Listener
}

func NewPProfServer(zl *zap.Logger) (*PProfServer, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", pprofPort))
	if err != nil {
		return nil, ErrFailedServerListenerInit(err)
	}

	psrv := &PProfServer{
		lis: lis,
		httpServer: &http.Server{
			Handler:           http.DefaultServeMux,
			ReadHeaderTimeout: 5 * time.Second,   //nolint:mnd
			ReadTimeout:       30 * time.Second,  //nolint:mnd
			WriteTimeout:      30 * time.Second,  //nolint:mnd
			IdleTimeout:       120 * time.Second, //nolint:mnd
		},
	}

	go func() {
		if err := psrv.httpServer.Serve(psrv.lis); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zl.Error("pprof server failed", zap.Error(err))
		}
	}()

	zl.Info("pprof server started", zap.Int("port", pprofPort))

	return psrv, nil
}

func (psrv *PProfServer) Close() error {
	var errs error

	if psrv.httpServer != nil {
		errs = multierr.Append(errs, psrv.httpServer.Close())
	}

	if psrv.lis != nil {
		if err := psrv.lis.Close(); !errors.Is(err, net.ErrClosed) {
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}
