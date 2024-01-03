package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

const keyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Printf("%s: Got root request \n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hi from root!")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Printf("%s: got the hello request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello, Http! Let's make it fun")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	ctx, cancelCtx := context.WithCancel(context.Background())
	serverOne := &http.Server{
		Addr:    ":8085",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	servertwo := &http.Server{
		Addr:    ":8086",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	serverthree := &http.Server{
		Addr:    ":8087",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := serverOne.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed \n")
		} else if err != nil {
			fmt.Printf("error starting server: %s \n", err)
			os.Exit(1)
		}
		cancelCtx()
	}()

	go func() {
		err := servertwo.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed \n")
		} else if err != nil {
			fmt.Printf("error starting server: %s \n", err)
			os.Exit(1)
		}
		cancelCtx()
	}()

	go func() {
		err := serverthree.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed \n")
		} else if err != nil {
			fmt.Printf("error starting server: %s \n", err)
			os.Exit(1)
		}
		cancelCtx()
	}()

	<-ctx.Done()
}
