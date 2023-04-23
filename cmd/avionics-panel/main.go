package main

import (
	"context"
	"net/http"
	"os"

	"github.com/wt-tools/avionics-panel/config"
	"github.com/wt-tools/wtscope/input/state"
	"github.com/wt-tools/wtscope/net/poll"

	"github.com/grafov/kiwi"
)

func main() {
	ctx := context.Background()
	kiwi.SinkTo(os.Stdout, kiwi.AsLogfmt()).Start()
	log := kiwi.New()
	conf := config.New()
	log.Log("status", "prepare avionics for start", "config", "xxx")
	errch := make(chan error, 8) // XXX разделить по компонентам
	defaultPolling := poll.New(http.DefaultClient, errch)
	stateWorker := state.New(conf, defaultPolling, errch)
	go defaultPolling.Do()
	go stateWorker.Grab(ctx)
	for {
	}
}
