package bonneFete

import "github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"

const (
	maxMessageSize = 1024
)

type WSHub struct {
	Subscribe    chan *WSClient
	Unsubscribe  chan *WSClient
	Clients      map[*WSClient]bool
	BroadcastMsg chan []byte
}

func newWSHub() *WSHub {
	return &WSHub{
		Subscribe:    make(chan *WSClient),
		Unsubscribe:  make(chan *WSClient),
		Clients:      make(map[*WSClient]bool),
		BroadcastMsg: make(chan []byte),
	}
}

func (h *WSHub) run() {
	for {
		select {
		case client := <-h.Subscribe:
			logger.Get().Info("Setting up new subscriber to /api/bonneFete/ws websocket")
			h.Clients[client] = true
			go client.writePump()
			go client.readPump()

		case client := <-h.Unsubscribe:
			if _, ok := h.Clients[client]; ok {
				logger.Get().Info("Unsubscribing client")
				delete(h.Clients, client)
			}

		case b := <-h.BroadcastMsg:
			for client := range h.Clients {
				client.msg <- b
			}
		}
	}
}
