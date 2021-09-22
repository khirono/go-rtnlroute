package rtnlroute

import (
	"net"
	"syscall"

	"github.com/khirono/go-nl"
)

type Request struct {
	Header Header
	Attrs  nl.AttrList
}

func (r *Request) AddDst(dst *net.IPNet) error {
	var dstip []byte
	switch len(dst.IP) {
	case net.IPv4len:
		r.Header.Family = syscall.AF_INET
		dstip = dst.IP.To4()
	case net.IPv6len:
		r.Header.Family = syscall.AF_INET6
		dstip = dst.IP.To16()
	}
	dstlen, _ := dst.Mask.Size()
	r.Header.Dstlen = uint8(dstlen)
	r.Attrs = append(r.Attrs, nl.Attr{
		Type:  syscall.RTA_DST,
		Value: nl.AttrBytes(dstip),
	})
	return nil
}

func (r *Request) AddIfName(ifname string) error {
	ifindex, err := nl.IfnameToIndex(ifname)
	if err != nil {
		return err
	}
	r.Attrs = append(r.Attrs, nl.Attr{
		Type:  syscall.RTA_OIF,
		Value: nl.AttrU32(ifindex),
	})
	return nil
}
