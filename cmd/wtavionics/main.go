package main

import (
	"context"
	"os"
	"time"

	"github.com/wt-tools/wtavionics/ui"
	"github.com/wt-tools/wtscope/config"
	"github.com/wt-tools/wtscope/input/indicators"
	"github.com/wt-tools/wtscope/input/state"
	"github.com/wt-tools/wtscope/net/poll"

	"github.com/grafov/kiwi"
)

func main() {
	ctx := context.Background()
	kiwi.SinkTo(os.Stdout, kiwi.AsLogfmt()).Start()
	l := kiwi.New()
	l.Log("status", "prepare avionics for start", "config", config.ConfPath)
	errch := make(chan error, 8) // XXX разделить по компонентам
	cfg, err := config.Load(errch)
	if err != nil {
		l.Log("status", "can't load configuration", "path", config.ConfPath)
		if err := config.CreateIfAbsent(); err != nil {
			os.Exit(1)
		}
		l.Log("status", "default configuration created")
		l.Log("hint", "check the file and fill it with your real config values", "path", config.ConfPath)
		os.Exit(0)
	}
	l.Log("status", "configuration loaded", "config", cfg.Dump())
	go showErrors(l, errch)
	defaultPolling := poll.New(poll.SetLogger(errch),
		poll.SetLoopDelay(250*time.Millisecond), poll.SetProblemDelay(4*time.Second))
	go defaultPolling.Do()
	gui := ui.Init(ctx, l)

	// TODO сделать сервисы входных данных отключаемыми в конфиге
	{
		stateSvc := state.New(cfg, defaultPolling, errch)
		go stateSvc.Grab(ctx)
		indSvc := indicators.New(cfg, defaultPolling, errch)
		go indSvc.Grab(ctx)
		gui.UpdateAvionics(ctx, stateSvc, indSvc)
	}
	gui.Run(ctx)
}

func showErrors(log *kiwi.Logger, errs chan error) {
	l := log.New()
	for {
		err := <-errs
		l.Log("problem", "library parser failed", "error", err)
	}
}
