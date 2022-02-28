package bonneFete

import (
	"fmt"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"
	"github.com/gorilla/websocket"
	"time"
)

var (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = time.Minute

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = 5 * time.Second
)

type WSClient struct {
	ws  *websocket.Conn
	hub *WSHub
	msg chan []byte
}

// writePump is the async flow to let the server communicate with the client.
// It sends pings and broadcast data to the client.
func (c *WSClient) writePump() {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		c.hub.Unsubscribe <- c
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.msg:
			fmt.Println("New message to broadcast!")
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.ws.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}

		case <-pingTicker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump is the async flow to let the client communicate with the server
// and ensures that the client answers "pong" the the server's "ping".
// If connection is broken, the client gets unsubscribed and connection is closed
func (c *WSClient) readPump() {
	defer func() {
		// when ReadLimit is reached, unsubscribe and close the connection
		c.hub.Unsubscribe <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)

	// reset read timeout upon receiving pong from the client
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Get().Error(err)
			}
			break
		}
		logger.Get().Info("Got client message. Client messages are not treated by the server.")
	}
}
