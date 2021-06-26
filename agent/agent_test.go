package agent

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var upgrader = websocket.Upgrader{}

func TestAgent(t *testing.T) {
	var c *websocket.Conn
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		c, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Log("upgrade:", err)
			return
		}
		defer c.Close()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			err := c.WriteMessage(websocket.TextMessage, []byte(time.Now().String()))
			if err != nil {
				t.Log("write:", err)
				break
			}
			<-ticker.C
		}
	})

	server := http.Server{
		Addr: "localhost:8080",
	}

	ln, err := net.Listen("tcp", server.Addr)
	require.NoError(t, err)

	errch := make(chan error)
	go func() {
		errch <- server.Serve(ln)
	}()

	agent := Agent{ServerAddr: "localhost:8080"}
	go agent.Poller(context.Background())

	server.RegisterOnShutdown(func() {
		// Cleanly close the connection by sending a close message and then
		// waiting (with timeout) for the server to close the connection.
		err := c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
		if err != nil {
			t.Log("write close:", err)
			return
		}
		return
	})

	time.Sleep(2 * time.Second)

	server.Shutdown(context.Background())

	assert.Equal(t, http.ErrServerClosed, <-errch)
}
