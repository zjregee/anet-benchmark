package runner

type Message struct {
	Message string
}

type Client interface {
	Echo(req *Message) (*Message, error)
}

type ClientNewer func(mode Mode, network, address string) Client

type Server interface {
	Run() error
}

type ServerNewer func(mode Mode, network, address string) Server

type EchoOnce func(req *Message) (*Message, error)

type Conn interface {
	Close() error
}

type Mode int

const (
	MODE_RPC Mode = 1
	MODE_MUX Mode = 2
)
