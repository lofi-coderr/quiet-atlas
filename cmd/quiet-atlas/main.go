package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"quietatlas/internal/server"
	"quietatlas/web"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "local HTTP listen address")
	flag.Parse()
	assets, err := web.Assets()
	if err != nil {
		log.Fatal(err)
	}
	index, err := assets.Open("index.html")
	if err != nil {
		log.Fatal("frontend missing: run cd web && npm ci && npm run build")
	}
	if err := index.Close(); err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{Addr: *addr, Handler: server.New(assets), ReadHeaderTimeout: 5 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdown); err != nil {
			log.Printf("shutdown: %v", err)
		}
	}()
	log.Printf("Quiet Atlas is ready at http://%s", *addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
