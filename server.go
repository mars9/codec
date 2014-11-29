package codec

import (
	"bufio"
	"encoding/binary"
	"io"
	"net/rpc"
	"sync"

	"github.com/golang/protobuf/proto"
)

const MaxVarint = binary.MaxVarintLen64

func writeMessage(w io.Writer, m interface{}) error {
	data, err := proto.Marshal(m.(proto.Message))
	if err != nil {
		return err
	}

	buf := make([]byte, MaxVarint)
	n := binary.PutUvarint(buf[:], uint64(len(data)))
	if _, err = w.Write(buf[:n]); err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

type messageReader interface {
	io.ByteReader
	io.Reader
}

func readMessage(r messageReader, m interface{}) error {
	size, err := binary.ReadUvarint(r)
	if err != nil {
		return err
	}

	data := make([]byte, size)
	if _, err = r.Read(data); err != nil {
		return err
	}
	if m != nil {
		return proto.Unmarshal(data, m.(proto.Message))
	}
	return nil
}

func NewServerCodec(rwc io.ReadWriteCloser) rpc.ServerCodec {
	return &serverCodec{
		r: bufio.NewReader(rwc),
		w: rwc,
		c: rwc,
	}
}

type serverCodec struct {
	r      messageReader
	c      io.Closer
	writer sync.Mutex
	w      io.Writer
}

func (s *serverCodec) ReadRequestHeader(req *rpc.Request) (err error) {
	var header RequestHeader
	if err = readMessage(s.r, &header); err != nil {
		return err
	}
	req.ServiceMethod = *header.Method
	req.Seq = *header.Seq
	return nil
}

func (s *serverCodec) ReadRequestBody(body interface{}) error {
	return readMessage(s.r, body)
}

func (s *serverCodec) WriteResponse(resp *rpc.Response, body interface{}) (err error) {
	header := ResponseHeader{
		Method: &resp.ServiceMethod,
		Seq:    &resp.Seq,
	}
	if resp.Error != "" {
		header.Error = &resp.Error
	}

	s.writer.Lock()
	if err = writeMessage(s.w, &header); err != nil {
		s.writer.Unlock()
		return err
	}
	err = writeMessage(s.w, body)
	s.writer.Unlock()
	return err
}

func (s *serverCodec) Close() error { return s.c.Close() }
