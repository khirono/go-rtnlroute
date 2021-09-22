package rtnlroute

import (
	"syscall"

	"github.com/khirono/go-nl"
)

func Create(c *nl.Client, r *Request) error {
	flags := syscall.NLM_F_CREATE
	flags |= syscall.NLM_F_EXCL
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(syscall.RTM_NEWROUTE, flags)
	err := req.Append(r.Header)
	if err != nil {
		return err
	}
	err = req.Append(r.Attrs)
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	return err
}

func Remove(c *nl.Client, r *Request) error {
	flags := syscall.NLM_F_ACK
	req := nl.NewRequest(syscall.RTM_DELROUTE, flags)
	err := req.Append(r.Header)
	if err != nil {
		return err
	}
	err = req.Append(r.Attrs)
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	return err
}

func GetAll(c *nl.Client, r *Request) ([]Route, error) {
	flags := syscall.NLM_F_DUMP
	req := nl.NewRequest(syscall.RTM_GETROUTE, flags)
	err := req.Append(r.Header)
	if err != nil {
		return nil, err
	}
	err = req.Append(r.Attrs)
	if err != nil {
		return nil, err
	}
	req.AppendReplyType(syscall.RTM_NEWROUTE)
	rsps, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	var rs []Route
	for _, rsp := range rsps {
		h, err := DecodeHeader(rsp.Body)
		if err != nil {
			return nil, err
		}
		r, err := DecodeRoute(rsp.Body[h.Len():])
		if err != nil {
			return nil, err
		}
		r.Header = h
		rs = append(rs, *r)
	}
	return rs, nil
}
