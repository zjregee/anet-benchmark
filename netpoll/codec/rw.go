package codec

import (
	"encoding/binary"

	"anet-benchmark/runner"
	"github.com/cloudwego/netpoll"
)

func Encode(writer netpoll.Writer, msg *runner.Message) error {
	header, err := writer.Malloc(4)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(header, uint32(len(msg.Message)))
	writer.WriteString(msg.Message)
	err = writer.Flush()
	return err
}

func Decode(reader netpoll.Reader, msg *runner.Message) error {
	bLen, err := reader.Next(4)
	if err != nil {
		return err
	}
	len := int(binary.BigEndian.Uint32(bLen))
	msg.Message, err = reader.ReadString(len)
	if err != nil {
		return err
	}
	err = reader.Release()
	return err
}
