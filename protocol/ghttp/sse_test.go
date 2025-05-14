package ghttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSSEGet(t *testing.T) {
	// 初始化 SSE 客户端
	inst := NewSSEInst(&SSEInstConfig{
		Module: "testModule",
		Host:   "127.0.0.1",
		Retry:  3,
	})

	counter := 0
	sseServer := createSSETestServer(
		10*time.Millisecond,
		func(w io.Writer) error {
			if counter == 100 {
				inst.Es().Close()
				return fmt.Errorf("stop sending events")
			}
			_, err := fmt.Fprintf(w, "id: %v\ndata: The time is %s\n\n", counter, time.Now().Format(time.UnixDate))
			counter++
			return err
		},
	)
	defer sseServer.Close()

	ctx := context.Background()
	err := inst.Es().
		SetURL(sseServer.URL).
		OnOpen(inst.NewOpenHandler(ctx)).
		OnError(inst.NewErrorHandler(ctx)).
		OnMessage(inst.NewMessageHandler(ctx), nil).Get()
	assert.Nil(t, err)
}

func createSSETestServer(ticker time.Duration, fn func(io.Writer) error) *httptest.Server {
	return createTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// for local testing allow it
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Create a channel for client disconnection
		clientGone := r.Context().Done()

		rc := http.NewResponseController(w)
		tick := time.NewTicker(ticker)
		defer tick.Stop()
		for {
			select {
			case <-clientGone:
				fmt.Println("Client disconnected")
				return
			case <-tick.C:
				if err := fn(w); err != nil {
					fmt.Println(err)
					return
				}
				if err := rc.Flush(); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	})
}

func createTestServer(fn func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(fn))
}
