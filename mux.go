package bootstrap

import (
	"context"
	"go.shu.run/bootstrap/logger"
	"go.shu.run/bootstrap/mux"
	"go.uber.org/fx"
	"net"
	"net/http"
)

type MuxConfig struct {
	ListenAt string `json:"listen_at"`
}

func StartMux(mux *mux.Mux, log logger.Logger, cfg MuxConfig, fc fx.Lifecycle) {
	ms := &MuxServer{
		log:     log.Prefix("Mux"),
		cfg:     cfg,
		handler: mux,
	}
	fc.Append(fx.Hook{
		OnStart: ms.OnStart,
		OnStop:  ms.OnStop,
	})
}

type MuxServer struct {
	fx.In
	log     logger.Logger
	cfg     MuxConfig
	handler *mux.Mux
	server  *http.Server
}

func (m *MuxServer) OnStart(ctx context.Context) error {
	m.server = &http.Server{
		Addr:    m.cfg.ListenAt,
		Handler: m.handler,
		BaseContext: func(ln net.Listener) context.Context {
			m.log.Debugf("connect: %s", ln.Addr().String())
			return ctx
		},
	}
	m.handler.SetLogger(m.log)
	m.log.Infof("http server starting...")
	go m.server.ListenAndServe()
	m.log.Infof("listen at: %s", m.cfg.ListenAt)
	return nil
}

func (m *MuxServer) OnStop(ctx context.Context) error {
	return m.server.Shutdown(ctx)
}