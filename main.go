package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"context"
	"log"

	"golang.org/x/sync/errgroup"
)

func main() {
	if len(os.Args) !=2 {
		fmt.Println("Please specify a port number.")
		os.Exit(1)
	}
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	if err := run(context.Background(),l); err != nil {
		log.Printf("failed to teminate server: %s", err)
		os.Exit(1)
	}
}

// runはサーバーを起動して、シグナルを待つ
func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		// Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	// 別のgoroutineでHTTPサーバーを起動する
	eg.Go(func() error {
		if err := s.Serve(l); err != nil &&
			err != http.ErrServerClosed {
				log.Printf("failed to listen and serve: %s", err)
				return err
		}
		return nil
	})

	// 他のgoroutineからのシグナル（終了通知）を待つ
	<-ctx.Done()
	// シグナルを受け取ったらサーバーを終了する
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown server: %s", err)
	}
	return eg.Wait()
}