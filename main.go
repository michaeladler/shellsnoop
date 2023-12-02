//go:build amd64

package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/michaeladler/shellsnoop/bpf"
	"github.com/michaeladler/shellsnoop/internal/log"
	"github.com/michaeladler/shellsnoop/internal/server"
	"github.com/michaeladler/shellsnoop/internal/storage"
	"github.com/mkideal/cli"
	"golang.org/x/sync/errgroup"
)

type argT struct {
	cli.Helper
	Uid    int    `cli:"u,uid" usage:"only snoop into processes owned by the given uid"`
	Socket string `cli:"s,socket" usage:"path for unix-domain socket to set up" dft:"/run/shellsnoop/shellsnoop.sock"`
	Level  string `cli:"l,log-level" usage:"set log level (possible values: debug, info, warn, error, fatal)" dft:"info"`
}

var (
	// provided by go linker
	Commit = ""
)

func main() {
	cli.Run(new(argT), func(cliCtx *cli.Context) error {

		argv := cliCtx.Argv().(*argT)

		var uid *int
		if cliCtx.IsSet("-u", "--uid") {
			uid = &argv.Uid
		}

		logger := log.Logger
		if err := log.SetLevel(argv.Level); err != nil {
			return err
		}

		logger.Info("Starting shellsnoop", "uid", argv.Uid, "socket", argv.Socket, "commit", Commit)

		if os.Geteuid() != 0 {
			return errors.New("ERROR: This program must be run as root user")
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		g, ctx := errgroup.WithContext(ctx)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		g.Go(func() error {
			if err := bpf.Run(ctx, uid); err != nil {
				return err
			}
			return nil
		})

		g.Go(func() error {
			return server.StartServer(ctx, argv.Socket, uid)
		})

		g.Go(func() error {
			running := true
			for running {
				select {
				case <-ctx.Done():
					running = false
				case <-time.After(10 * time.Minute):
					storage.GC()
				}
			}
			return nil
		})

		go func() {
			_ = g.Wait()
			sigChan <- syscall.SIGTERM
		}()

		<-sigChan

		// we are done, cancel our context to stop the goroutines
		cancel()

		if err := g.Wait(); err != nil {
			logger.Error("Something went wrong", "err", err)
		}
		return nil
	})
}
