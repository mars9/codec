package codec

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/golang/protobuf/proto"
)

const defaultBufferSize = 4 * 1024

// A Decoder manages the receipt of type and data information read from the
// remote side of a connection.
type Decoder struct {
	r    *bufio.Reader
	v    uint64
	size int
	buf  []byte
}

// NewDecoder returns a new decoder that reads from the io.Reader. It the
// argument io.Reader is a *bufio.Reader with large enough size, it uses the
// underlying *bufio.Reader. Otherwise it creates a *bufio.Reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReaderSize(r, defaultBufferSize)}
}

// Decode reads the next value from the input stream and stores it in the
// data represented by the empty interface value. If m is nil, the value
// will be discarded. Otherwise, the value underlying m must be a pointer
// to the correct type for the next data item received.
func (d *Decoder) Decode(m proto.Message) (err error) {
	if d.buf, err = readFull(d.r, d.buf); err != nil {
		return err
	}
	return proto.Unmarshal(d.buf, m)
}

func readFull(r *bufio.Reader, data []byte) ([]byte, error) {
	var err error
	v, err := binary.ReadUvarint(r)
	if err != nil {
		return data, err
	}
	size := int(v)
	if len(data) < size {
		data = make([]byte, size)
	}
	data = data[:size]

	_, err = io.ReadFull(r, data)
	return data, err
}

// An Encoder manages the transmission of type and data information to the
// other side of a connection.
type Encoder struct {
	size [binary.MaxVarintLen64]byte
	buf  *proto.Buffer
	w    io.Writer
}

// NewEncoder returns a new encoder that will transmit on the io.Writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{buf: &proto.Buffer{}, w: w}
}

// Encode transmits the data item represented by the empty interface value,
// guaranteeing that all necessary type information has been transmitted
// first.
func (e *Encoder) Encode(m proto.Message) (err error) {
	if err = e.buf.Marshal(m); err != nil {
		e.buf.Reset()
		return err
	}
	err = e.writeFrame(e.buf.Bytes())
	e.buf.Reset()
	return err
}

func (e *Encoder) writeFrame(data []byte) (err error) {
	n := binary.PutUvarint(e.size[:], uint64(len(data)))
	if _, err = e.w.Write(e.size[:n]); err != nil {
		return err
	}
	_, err = e.w.Write(data)
	return err
}
