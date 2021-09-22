package rtnlroute

import (
	"net"
	"syscall"

	"github.com/khirono/go-nl"
)

const (
	RTA_PREF = 20
)

// TODO:
type Route struct {
	Header  Header
	Table   uint32
	Pref    uint8
	Dst     net.IP
	Gw      net.IP
	Prefsrc net.IP
	Oif     uint32
	Prio    uint32
	Cache   *Cacheinfo
}

func DecodeRoute(b []byte) (*Route, error) {
	var r Route
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return nil, err
		}
		switch hdr.MaskedType() {
		case syscall.RTA_DST:
			r.Dst = net.IP(b[n : n+int(hdr.Len-4)])
		case syscall.RTA_OIF:
			r.Oif = native.Uint32(b[n:])
		case syscall.RTA_GATEWAY:
			r.Gw = net.IP(b[n : n+int(hdr.Len-4)])
		case syscall.RTA_PRIORITY:
			r.Prio = native.Uint32(b[n:])
		case syscall.RTA_PREFSRC:
			r.Prefsrc = net.IP(b[n : n+int(hdr.Len-4)])
		case syscall.RTA_CACHEINFO:
			c, err := DecodeCacheinfo(b[n:])
			if err != nil {
				break
			}
			r.Cache = c
		case syscall.RTA_TABLE:
			r.Table = native.Uint32(b[n:])
		case RTA_PREF:
			r.Pref = b[n]
		}
		b = b[hdr.Len.Align():]
	}
	return &r, nil
}

type Cacheinfo struct {
	Clntref uint32
	Lastuse uint32
	Expires uint32
	Error   uint32
	Used    uint32
	ID      uint32
	Ts      uint32
	Tsage   uint32
}

func DecodeCacheinfo(b []byte) (*Cacheinfo, error) {
	var c Cacheinfo
	c.Clntref = native.Uint32(b[0:4])
	c.Lastuse = native.Uint32(b[4:8])
	c.Expires = native.Uint32(b[8:12])
	c.Error = native.Uint32(b[12:16])
	c.Used = native.Uint32(b[16:20])
	c.ID = native.Uint32(b[20:24])
	c.Ts = native.Uint32(b[24:28])
	c.Tsage = native.Uint32(b[28:32])
	return &c, nil
}
