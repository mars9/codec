package codec

import (
	"bufio"
	"io"
	"net/rpc"
	"sync"
)

type clientCodec struct {
	r      messageReader
	c      io.Closer
	writer sync.Mutex
	w      io.Writer
}

func NewClientCodec(rwc io.ReadWriteCloser) rpc.ClientCodec {
	return &clientCodec{
		r: bufio.NewReader(rwc),
		w: rwc,
		c: rwc,
	}
}

func (c *clientCodec) WriteRequest(req *rpc.Request, body interface{}) (err error) {
	header := RequestHeader{
		Method: &req.ServiceMethod,
		Seq:    &req.Seq,
	}

	c.writer.Lock()
	if err = writeMessage(c.w, &header); err != nil {
		c.writer.Unlock()
		return err
	}
	err = writeMessage(c.w, body)
	c.writer.Unlock()
	return err
}

func (c *clientCodec) ReadResponseHeader(resp *rpc.Response) (err error) {
	var header ResponseHeader
	if err = readMessage(c.r, &header); err != nil {
		return err
	}
	resp.ServiceMethod = *header.Method
	resp.Seq = *header.Seq
	if header.Error != nil {
		resp.Error = *header.Error
	}
	return err
}

func (c *clientCodec) ReadResponseBody(body interface{}) error {
	return readMessage(c.r, body)
}

func (c *clientCodec) Close() error { return c.c.Close() }
