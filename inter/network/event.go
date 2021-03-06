package network

import "net"

// SocketEvent is the generic struct for events
// by this socket
//
// Current events:
//		Name				-> Data-Type
// 		close 				-> nil
//		error				-> error
//		newClient			-> *Client
//		client.close		-> [0: *client, 1:nil]
//		client.error		-> {*client, error}
// 		client.command		-> {*client, *Command}
//		client.command.*	-> {*client, *Command}
//		client.data			-> {*client, string}
type SocketEvent struct {
	Name string
	Data interface{}
}

type SocketUDPEvent struct {
	Name string
	Addr *net.UDPAddr
	Data interface{}
}

// ClientEvent is the generic struct for events
// by this Client
type ClientEvent struct {
	Name string
	Data interface{}
}

type EventError struct {
	Error error
}
type EventNewClient struct {
	Client *Client
}

type EventClientClose struct {
	Client *Client
}
type EventClientError struct {
	Client *Client
	Error  error
}
type EventClientCommand struct {
	Client  *Client
	Command *CommandFESL // If TLS (theater then we ignore payloadID - it is always 0x0)
}

type EventClientData struct {
	Client *Client
	Data   string
}

func (c *Client) FireClientClose(event ClientEvent) SocketEvent {
	return SocketEvent{
		Name: "client.close",
		Data: EventClientClose{Client: c},
	}
}

func (c *Client) FireClose() ClientEvent {
	return ClientEvent{
		Name: "close",
		Data: c,
	}
}

func (c *Client) FireError(err error) ClientEvent {
	return ClientEvent{
		Name: "error",
		Data: err,
	}
}

func (c *Client) FireClientData(event ClientEvent) SocketEvent {
	return SocketEvent{
		Name: "client.data",
		Data: EventClientData{
			Client: c,
			Data:   event.Data.(string),
		},
	}
}

func (c *Client) FireClientCommand(event ClientEvent) SocketEvent {
	return SocketEvent{
		Name: "client." + event.Name,
		Data: EventClientCommand{
			Client:  c,
			Command: event.Data.(*CommandFESL),
		},
	}
}

func (c *Client) FireSomething(event ClientEvent) SocketEvent {
	var interfaceSlice = make([]interface{}, 2)
	interfaceSlice[0] = c
	interfaceSlice[1] = event.Data
	return SocketEvent{
		Name: "client." + event.Name,
		Data: interfaceSlice,
	}
}

// ===

func (s *Socket) FireError(err error) SocketEvent {
	return SocketEvent{
		Name: "error",
		Data: EventError{Error: err},
	}
}

func (s *Socket) FireClose() SocketEvent {
	return SocketEvent{
		Name: "close",
		Data: nil,
	}
}

func (s *Socket) FireNewClient(client *Client) SocketEvent {
	return SocketEvent{
		Name: "newClient",
		Data: EventNewClient{Client: client},
	}
}
