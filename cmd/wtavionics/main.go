package main

import (
	"context"
	"net/http"
	"os"

	"github.com/wt-tools/wtavionics/config"
	"github.com/wt-tools/wtavionics/ui"
	"github.com/wt-tools/wtscope/input/indicators"
	"github.com/wt-tools/wtscope/input/state"
	"github.com/wt-tools/wtscope/net/poll"

	"github.com/grafov/kiwi"
)

func main() {
	ctx := context.Background()
	kiwi.SinkTo(os.Stdout, kiwi.AsLogfmt()).Start()
	l := kiwi.New()
	conf := config.New()
	l.Log("status", "prepare avionics for start", "config", "xxx")
	errch := make(chan error, 8) // XXX разделить по компонентам
	defaultPolling := poll.New(http.DefaultClient, errch, 2, 2)
	stateSvc := state.New(conf, defaultPolling, errch)
	indSvc := indicators.New(conf, defaultPolling, errch)
	go defaultPolling.Do()
	go stateSvc.Grab(ctx)
	go indSvc.Grab(ctx)
	gui := ui.Init(ctx, l)
	gui.UpdateAvionics(ctx, stateSvc, indSvc)
	gui.Run(ctx)
}
