package codec

import proto "code.google.com/p/goprotobuf/proto"

type RequestHeader struct {
	Method           *string `protobuf:"bytes,1,req,name=method" json:"method,omitempty"`
	Seq              *uint64 `protobuf:"varint,2,req,name=seq" json:"seq,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *RequestHeader) Reset()         { *m = RequestHeader{} }
func (m *RequestHeader) String() string { return proto.CompactTextString(m) }
func (*RequestHeader) ProtoMessage()    {}

type ResponseHeader struct {
	Method           *string `protobuf:"bytes,1,req,name=method" json:"method,omitempty"`
	Seq              *uint64 `protobuf:"varint,2,req,name=seq" json:"seq,omitempty"`
	Error            *string `protobuf:"bytes,3,opt,name=error" json:"error,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ResponseHeader) Reset()         { *m = ResponseHeader{} }
func (m *ResponseHeader) String() string { return proto.CompactTextString(m) }
func (*ResponseHeader) ProtoMessage()    {}
