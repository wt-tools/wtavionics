package main

import (
	"context"
	"os"
	"time"

	"github.com/wt-tools/wtavionics/config"
	"github.com/wt-tools/wtavionics/ui"
	"github.com/wt-tools/wtscope/input/hudmsg"
	"github.com/wt-tools/wtscope/input/indicators"
	"github.com/wt-tools/wtscope/input/state"
	"github.com/wt-tools/wtscope/net/dedup"
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
	go showErrors(l, errch)
	defaultPolling := poll.New(poll.SetLogger(errch),
		poll.SetLoopDelay(250*time.Millisecond), poll.SetProblemDelay(4*time.Second))
	go defaultPolling.Do()
	gui := ui.Init(ctx, l)

	// TODO сделать сервисы входных данных отключаемыми в конфиге
	{
		stateSvc := state.New(conf, defaultPolling, errch)
		go stateSvc.Grab(ctx)
		indSvc := indicators.New(conf, defaultPolling, errch)
		go indSvc.Grab(ctx)
		gui.UpdateAvionics(ctx, stateSvc, indSvc)
	}

	{
		battleSvc := hudmsg.New(conf, defaultPolling, dedup.New(), errch)
		go battleSvc.Grab(ctx)
		gui.UpdateBattleLog(ctx, battleSvc)
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
