package perf_test

import (
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/jasonkwh/perf"
	"go.uber.org/zap"
)

func TestNewPProfServer(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "listener already bound",
			wantErr: true,
		},
		{
			name:    "success",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := assert.New(t)

			if tt.wantErr {
				lis, listenErr := net.Listen("tcp", ":6060")
				if !as.NoError(listenErr) {
					return
				}

				t.Cleanup(func() { lis.Close() })
			}

			srv, err := perf.NewPProfServer(zap.NewNop())
			if tt.wantErr {
				as.Nil(srv)
				as.ErrorContains(err, "failed to initialise server listener")

				return
			}

			as.NoError(err)
			if !as.NotNil(srv) {
				return
			}

			t.Cleanup(func() { as.NoError(srv.Close()) })

			client := &http.Client{Timeout: time.Second}
			deadline := time.Now().Add(2 * time.Second)

			var (
				resp   *http.Response
				getErr error
			)

			for {
				resp, getErr = client.Get("http://127.0.0.1:6060/debug/pprof/")
				if getErr == nil || time.Now().After(deadline) {
					break
				}

				time.Sleep(5 * time.Millisecond)
			}

			if !as.NoError(getErr) {
				return
			}

			t.Cleanup(func() { resp.Body.Close() })
			as.Equal(http.StatusOK, resp.StatusCode)
		})
	}
}
