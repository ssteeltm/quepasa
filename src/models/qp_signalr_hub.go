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

// not used
func SignalRHubFactory() signalr.HubInterface {
	return MessagesSignalRHub
}

func (source *QpSignalRHub) IsInterfaceNil() bool {
	return source == nil
}

func (source *QpSignalRHub) OnConnected(ConnectionId string) {
	info, _ := source.Logger()
	info.Log("connection", ConnectionId, "status", "connected")

	Proxy[ConnectionId] = source.Clients().Caller()
}

func (source *QpSignalRHub) OnDisconnected(ConnectionId string) {
	info, _ := source.Logger()
	info.Log("connection", ConnectionId, "status", "disconnected")

	delete(Tokens, ConnectionId)
	delete(Proxy, ConnectionId)
}

func (source *QpSignalRHub) TrySend(target string, args ...interface{}) {
	if source == nil {
		return
	}

	ConnectionId := source.ConnectionID()
	TrySend(ConnectionId, target, args)
}

func TrySend(ConnectionId string, target string, args ...interface{}) {
	proxy := Proxy[ConnectionId]
	if proxy != nil {
		proxy.Send(target, args)
	}
}

func (source *QpSignalRHub) GetToken() string {
	ConnectionId := source.ConnectionID()
	token := Tokens[ConnectionId]

	message := fmt.Sprintf("connection id: %s, token: %s", ConnectionId, token)
	TrySend(ConnectionId, "system", message)
	return token
}

func (source *QpSignalRHub) Token(token string) {
	ConnectionId := source.ConnectionID()
	Tokens[ConnectionId] = token

	info, _ := source.Logger()
	info.Log("connection", ConnectionId, "token", token)
}

func SignalRDispatch(token string, payload *whatsapp.WhatsappMessage) {
	for ConnectionId, _token := range Tokens {
		if _token != token {
			continue
		}

		TrySend(ConnectionId, "message", payload)
	}
}
