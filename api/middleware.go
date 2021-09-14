package api

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

// LogInit will hold the writer and will be used while writing logs.
type LogInit struct {
	Logpath io.Writer
}

// TimeoutMiddleware holds the http handler for which timeout has to be set.
type TimeoutMiddleware struct {
	Next http.Handler
}

const (
	timeouts = 30 * time.Second
)

func (tm TimeoutMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if tm.Next == nil {
		tm.Next = http.DefaultServeMux
	}

	ctx := r.Context()
	ctx, _ = context.WithTimeout(ctx, timeouts)
	r.WithContext(ctx)
	ch := make(chan struct{})
	go func() {
		tm.Next.ServeHTTP(w, r)
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		return
	case <-ctx.Done():
		w.WriteHeader(http.StatusRequestTimeout)
	}
}

// TimeoutHandler will be responsible for timing out the request if it takes annoying time to serve.
func TimeoutHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tCtx, tCancel := context.WithTimeout(ctx, timeouts)
		cCtx, cCancel := context.WithCancel(ctx)
		r.WithContext(ctx)

		defer tCancel()

		go func() {
			next.ServeHTTP(w, r)
			cCancel()
		}()
		select {
		case <-cCtx.Done():
			return
		case <-tCtx.Done():
			if err := tCtx.Err(); err == context.DeadlineExceeded {
				cCancel()
				w.WriteHeader(http.StatusGatewayTimeout)
				if _, err := w.Write([]byte("Request Timeout, it is taking more than anticipated time to serve request\n")); err != nil {
					log.Fatal(err)
				}
			}
		}
	})
}

// Logger will log all the request served by the neuron through API.
func (loger *LogInit) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		newlog := log.New(loger.Logpath, " [INFO ]", 0)
		newlog.SetPrefix(logformatter(r, start))
		newlog.Println()
	})
}

func logformatter(r *http.Request, start time.Time) string {
	return "[" + start.Format("2006-01-02 15:04:05") + "]" + " - " + r.Method + " " + r.RequestURI +
		" " + r.RemoteAddr + " " + time.Since(start).String() + " " + r.Proto
}
