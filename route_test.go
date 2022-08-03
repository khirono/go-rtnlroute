package rtnlroute

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"syscall"
	"testing"

	"github.com/khirono/go-nl"
	"github.com/khirono/go-rtnllink"
)

func Setup(t *testing.T, c *nl.Client) func() {
	peer := nl.Attr{
		Type: rtnllink.VETH_INFO_PEER | syscall.NLA_F_NESTED,
		Value: nl.Encoders{
			rtnllink.IfInfomsg{},
			&nl.Attr{
				Type:  syscall.IFLA_IFNAME,
				Value: nl.AttrString("bar"),
			},
		},
	}
	linkinfo := &nl.Attr{
		Type: syscall.IFLA_LINKINFO,
		Value: nl.AttrList{
			{
				Type:  rtnllink.IFLA_INFO_KIND,
				Value: nl.AttrString("veth"),
			},
			{
				Type: rtnllink.IFLA_INFO_DATA,
				Value: nl.AttrList{
					peer,
				},
			},
		},
	}
	err := rtnllink.Create(c, "foo", linkinfo)
	if err != nil {
		t.Fatal(err)
	}

	err = rtnllink.Up(c, "foo")
	if err != nil {
		t.Fatal(err)
	}

	return func() {
		err = rtnllink.Remove(c, "foo")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestStaticRoute(t *testing.T) {
	var wg sync.WaitGroup
	mux, err := nl.NewMux()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		mux.Close()
		wg.Wait()
	}()
	wg.Add(1)
	go func() {
		mux.Serve()
		wg.Done()
	}()

	conn, err := nl.Open(syscall.NETLINK_ROUTE)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := nl.NewClient(conn, mux)

	teardown := Setup(t, c)
	defer teardown()

	t.Run("ip route add 60.60.0.0/16 dev foo", func(t *testing.T) {
		r := &Request{
			Header: Header{
				Table:    syscall.RT_TABLE_MAIN,
				Scope:    syscall.RT_SCOPE_UNIVERSE,
				Protocol: syscall.RTPROT_STATIC,
				Type:     syscall.RTN_UNICAST,
			},
		}
		_, dst, err := net.ParseCIDR("60.60.0.0/16")
		if err != nil {
			t.Fatal(err)
		}
		err = r.AddDst(dst)
		if err != nil {
			t.Fatal(err)
		}
		err = r.AddIfName("foo")
		if err != nil {
			t.Fatal(err)
		}
		err = Create(c, r)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("ip route del 60.60.0.0/16 dev foo", func(t *testing.T) {
		r := &Request{
			Header: Header{
				Table:    syscall.RT_TABLE_MAIN,
				Scope:    syscall.RT_SCOPE_UNIVERSE,
				Protocol: syscall.RTPROT_STATIC,
				Type:     syscall.RTN_UNICAST,
			},
		}
		_, dst, err := net.ParseCIDR("60.60.0.0/16")
		if err != nil {
			t.Fatal(err)
		}
		err = r.AddDst(dst)
		if err != nil {
			t.Fatal(err)
		}
		err = r.AddIfName("foo")
		if err != nil {
			t.Fatal(err)
		}
		err = Remove(c, r)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestGetAll(t *testing.T) {
	var wg sync.WaitGroup
	mux, err := nl.NewMux()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		mux.Close()
		wg.Wait()
	}()
	wg.Add(1)
	go func() {
		mux.Serve()
		wg.Done()
	}()

	conn, err := nl.Open(syscall.NETLINK_ROUTE)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := nl.NewClient(conn, mux)

	rs, err := GetAll(c, &Request{})
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.MarshalIndent(rs, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("routes: %s\n", j)
}
