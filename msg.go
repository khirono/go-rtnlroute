package rtnlroute

const (
	SizeofHeader = 12
)

type Header struct {
	Family   uint8
	Dstlen   uint8
	Srclen   uint8
	Tos      uint8
	Table    uint8
	Protocol uint8
	Scope    uint8
	Type     uint8
	Flags    uint32
}

func DecodeHeader(b []byte) (Header, error) {
	var h Header
	h.Family = b[0]
	h.Dstlen = b[1]
	h.Srclen = b[2]
	h.Tos = b[3]
	h.Table = b[4]
	h.Protocol = b[5]
	h.Scope = b[6]
	h.Type = b[7]
	h.Flags = native.Uint32(b[8:12])
	return h, nil
}

func (h Header) Len() int {
	return SizeofHeader
}

func (h Header) Encode(b []byte) (int, error) {
	b[0] = h.Family
	b[1] = h.Dstlen
	b[2] = h.Srclen
	b[3] = h.Tos
	b[4] = h.Table
	b[5] = h.Protocol
	b[6] = h.Scope
	b[7] = h.Type
	native.PutUint32(b[8:12], h.Flags)
	return h.Len(), nil
}
