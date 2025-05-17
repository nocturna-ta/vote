package api

import (
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/router"
	"github.com/nocturna-ta/vote/config"
	"github.com/nocturna-ta/vote/internal/handler/api/controller"
	"github.com/nocturna-ta/vote/internal/usecases"
)

type Options struct {
	Cfg    config.MainConfig
	VoteUc usecases.VoteUseCases
}

type Handler struct {
	opts        *Options
	listenErrCh chan error
	myRouter    *router.FastRouter
}

func New(opts *Options) *Handler {
	handler := &Handler{
		opts: opts,
	}
	handler.myRouter = controller.New(&controller.Options{
		Prefix:         opts.Cfg.API.BasePath,
		Port:           opts.Cfg.Server.Port,
		ReadTimeout:    opts.Cfg.Server.ReadTimeout,
		WriteTimeout:   opts.Cfg.Server.WriteTimeout,
		RequestTimeout: opts.Cfg.API.APITimeout,
		EnableSwagger:  opts.Cfg.API.EnableSwagger,
		VoteUc:         opts.VoteUc,
	}).RegisterRoute()
	return handler
}
func (h *Handler) Run() {
	log.Infof("API Listening on %d", h.opts.Cfg.Server.Port)
	h.listenErrCh <- h.myRouter.StartServe()
}

func (h *Handler) ListenError() <-chan error {
	return h.listenErrCh
}
