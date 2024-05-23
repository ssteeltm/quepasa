package models

import (
	"fmt"

	whatsapp "github.com/nocodeleaks/quepasa/whatsapp"
	signalr "github.com/philippseith/signalr"
)

type QpSignalRHub struct {
	signalr.Hub
}

var MessagesSignalRHub = &QpSignalRHub{}
var Tokens = map[string]string{}
var Proxy = map[string]signalr.ClientProxy{}

func SignalRHubFactory() signalr.HubInterface {
	return MessagesSignalRHub
}

func (c *QpSignalRHub) IsInterfaceNil() bool {
	return c == nil
}

func ForEach(key any, value any) bool {
	fmt.Printf("item key: %s, value: %s\n", key, value)
	return true
}

func (c *QpSignalRHub) OnConnected(ConnectionId string) {
	fmt.Printf("%s connected\n", ConnectionId)
	Proxy[ConnectionId] = c.Clients().Caller()
}

func (c *QpSignalRHub) OnDisconnected(ConnectionId string) {
	fmt.Printf("%s disconnected\n", ConnectionId)
	delete(Tokens, ConnectionId)
	delete(Proxy, ConnectionId)
}

func (h *QpSignalRHub) TrySend(target string, args ...interface{}) {
	if h == nil {
		return
	}

	ConnectionId := h.ConnectionID()
	TrySend(ConnectionId, target, args)
}

func TrySend(ConnectionId string, target string, args ...interface{}) {
	proxy := Proxy[ConnectionId]
	if proxy != nil {
		proxy.Send(target, args...)
	}
}

func (h *QpSignalRHub) GetToken() string {
	ConnectionId := h.ConnectionID()
	token := Tokens[ConnectionId]

	message := fmt.Sprintf("connection id: %s, token: %s", ConnectionId, token)
	TrySend(ConnectionId, "system", message)
	return token
}

func (h *QpSignalRHub) Token(token string) {
	ConnectionId := h.ConnectionID()
	fmt.Printf("token: %s, for connection id: %s", token, ConnectionId)

	Tokens[ConnectionId] = token
}

func SignalRDispatch(token string, payload *whatsapp.WhatsappMessage) {
	for ConnectionId, _token := range Tokens {
		if _token != token {
			continue
		}

		TrySend(ConnectionId, "message", payload)
	}
}
