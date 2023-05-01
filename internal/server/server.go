package server

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/michaeladler/shellsnoop/internal/log"
	"github.com/michaeladler/shellsnoop/internal/storage"
	"golang.org/x/sync/errgroup"
)

func StartServer(ctx context.Context, socketFile string, uid *int) error {
	logger := log.Logger

	dir := filepath.Dir(socketFile)
	logger.Debug("Creating directory", "dir", dir)
	if err := os.Mkdir(dir, os.FileMode(0755)); err != nil && !os.IsExist(err) {
		logger.Warn("Failed to create directory", "err", err)
	}

	listener, err := net.Listen("unix", socketFile)
	if err != nil {
		return err
	}

	if uid != nil {
		logger.Debug("Changing owner", uid, *uid)
		_ = os.Chown(socketFile, *uid, 0)
		_ = os.Chmod(socketFile, 0600)
	}

	defer listener.Close()

	logger.Info("Listening on Unix domain socket", "socketFile", socketFile)

	g, ctx := errgroup.WithContext(ctx)

	running := true
	g.Go(func() error {
		// Accept and handle incoming connections.
		for running {
			conn, err := listener.Accept()
			if err != nil {
				if err != net.ErrClosed {
					logger.Error("Failed to accept connection", "err", err)
				}
				continue
			}
			defer conn.Close()

			// Handle the connection in a new goroutine.
			g.Go(func() error {
				err := handleConnection(conn)
				if err != nil {
					logger.Error("Something went wrong with the client", "err", err)
				}
				return nil
			})
		}
		return nil
	})

	doneChan := make(chan any)

	go func() {
		_ = g.Wait()
		doneChan <- nil
	}()

	// wait for quit event
	select {
	case <-doneChan:
		running = false
	case <-ctx.Done():
		running = false
	}

	listener.Close()
	_ = g.Wait()

	return nil
}

func handleConnection(conn net.Conn) error {
	// Read incoming data.
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return err
	}

	data := string(buffer[:n])

	pid, err := strconv.Atoi(data)
	if err != nil {
		return err
	}

	logger := log.Logger
	cmdline := storage.Get(pid)
	if cmdline == "" {
		logger.Warn("Last command not available", "pid", pid)
		_, _ = conn.Write([]byte("\x00"))
		return nil
	}

	logger.Info("Providing last command to client", "pid", pid)
	_, err = conn.Write([]byte(cmdline))
	if err != nil {
		return err
	}
	return nil
}
