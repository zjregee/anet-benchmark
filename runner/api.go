package runner

type Message struct {
	Message string
}

type Client interface {
	Echo(req *Message) (*Message, error)
}

type ClientNewer func(network, address string) Client

type Server interface {
	Run() error
}

type ServerNewer func(network, address string) Server

type EchoOnce func(req *Message) (*Message, error)

type Conn interface {
	Close() error
}
