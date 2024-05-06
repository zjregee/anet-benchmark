package codec

import (
	"io"
	"net"
	"sync"
	"bufio"
	"unsafe"
	"encoding/binary"

	"anet-benchmark/runner"
)

var connerpool sync.Pool

type Conner struct {
	header []byte
	rw     *bufio.ReadWriter
}

func NewConner(conn net.Conn) *Conner {
	if v := connerpool.Get(); v != nil {
		conner := v.(*Conner)
		conner.rw.Reader.Reset(conn)
		conner.rw.Writer.Reset(conn)
	}
	conner := &Conner{}
	conner.header = make([]byte, 4)
	conner.rw = bufio.NewReadWriter(bufio.NewReaderSize(conn, 8192), bufio.NewWriterSize(conn, 8192))
	return conner
}

func PutConner(conner *Conner) {
	conner.rw.Reader.Reset(nil)
	conner.rw.Writer.Reset(nil)
	connerpool.Put(conner)
}

func (c *Conner) Encode(msg *runner.Message) error {
	binary.BigEndian.PutUint32(c.header, uint32(len(msg.Message)))
	_, err := c.rw.Write(c.header)
	if err != nil {
		return err
	}
	_, err = c.rw.WriteString(msg.Message)
	if err != nil {
		return err
	}
	err = c.rw.Flush()
	return err
}

func (c *Conner) Decode(msg *runner.Message) error {
	bLen, err := c.rw.Peek(4)
	if err != nil {
		return err
	}
	_, err = c.rw.Discard(4)
	if err != nil {
		return err
	}
	len := binary.BigEndian.Uint32(bLen)
	payload := make([]byte, len)
	_, err = io.ReadFull(c.rw, payload)
	if err != nil {
		return err
	}
	msg.Message = b2s(payload)
	return nil
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
