package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

// implements http.Handler interface
type handler struct {
	name   string
	server *http.Server
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cmd := r.URL.Path[1:]
	fmt.Fprintf(w, "Hi there, handler=%s, req=%s\n", h.name, r.URL.Path[1:])
	if cmd == "close" {
		fmt.Printf("closing server %v by API\n", h.name)
		timer := time.NewTimer(0 * time.Second)
		go func() {
			// defered close, make sure this request can response successfuly.
			<-timer.C
			h.server.Close()
		}()
	}
}

func main() {
	sigsCh := make(chan os.Signal, 1)
	signal.Notify(sigsCh, syscall.SIGINT, syscall.SIGTERM)

	server1 := &http.Server{
		Addr: ":8080",
	}
	server1.Handler = &handler{name: "h1", server: server1}

	server2 := &http.Server{
		Addr: ":8082",
	}
	server2.Handler = &handler{name: "h2", server: server2}

	g, ctx := errgroup.WithContext(context.TODO())

	g.Go(func() error {
		go func() {
			<-ctx.Done()
			fmt.Printf("closing server 1 by graceful shutdown\n")
			server1.Close()
		}()
		return server1.ListenAndServe()
	})
	g.Go(func() error {
		go func() {
			<-ctx.Done()
			fmt.Printf("closing server 2 by graceful shutdown\n")
			server2.Close()
		}()
		return server2.ListenAndServe()
	})
	g.Go(func() error {
		select {
		case sig := <-sigsCh:
			fmt.Printf("signal recrived: %s\n", sig.String())
			return errors.New("signal")
		case <-ctx.Done():
			fmt.Printf("signal go routine exit\n")
			return nil
		}
	})

	err := g.Wait()

	fmt.Printf("exit reason: err=%v\n", err)
}
