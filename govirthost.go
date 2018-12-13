package controller

import (
	"errors"
	"fmt"
	"github.com/hunkeelin/govirt/govirtlib"
	"github.com/hunkeelin/klinutils"
	"github.com/hunkeelin/mtls/klinreq"
	"io/ioutil"
)

func (c *Conn) statevm(state, vm, host string) error {
	p := &govirtlib.PostPayload{
		Action: state,
		Domain: vm,
	}
	i := &klinreq.ReqInfo{
		Dest:       host,
		Dport:      klinutils.Stringtoport("govirthost"),
		Method:     "POST",
		Payload:    p,
		TrustBytes: c.tb,
		CertBytes:  c.cb,
		KeyBytes:   c.kb,
	}
	resp, err := klinreq.SendPayload(i)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(string(body))
		return errors.New("Failed, check logs on the govirthost server")
	}
	return nil
}
func (c *Conn) migrate(ori, dest, vm string) error {
	p := &govirtlib.PostPayload{
		Action: "Migrate",
		Target: dest,
		Domain: vm,
	}
	i := &klinreq.ReqInfo{
		Dest:       ori,
		Dport:      klinutils.Stringtoport("govirthost"),
		Method:     "POST",
		Payload:    p,
		TrustBytes: c.tb,
		CertBytes:  c.cb,
		KeyBytes:   c.kb,
	}
	resp, err := klinreq.SendPayload(i)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(string(body))
		return errors.New("Failed, check logs on the govirthost server")
	}
	return nil
}
func (c *Conn) Define(xml []byte, dest string) error {
	p := govirtlib.PostPayload{
		Action: "Define",
		Xml:    xml,
	}
	i := &klinreq.ReqInfo{
		Dest:       dest,
		Dport:      klinutils.Stringtoport("govirthost"),
		Method:     "POST",
		Payload:    p,
		TrustBytes: c.tb,
		CertBytes:  c.cb,
		KeyBytes:   c.kb,
	}
	resp, err := klinreq.SendPayload(i)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(string(body))
		return errors.New("Failed, check logs on the govirthost server")
	}
	return nil
}
