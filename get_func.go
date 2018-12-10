package controller

import (
	"errors"
	"fmt"
	"github.com/hunkeelin/govirt/govirtlib"
	"github.com/hunkeelin/klinutils"
	"github.com/hunkeelin/mtls/klinreq"
	"io/ioutil"
)

func (c *Conn) getxml(vm, host string) ([]byte, error) {
	var r []byte
	p := &govirtlib.GetPayload{
		Target: "xml",
		Domain: vm,
	}
	i := &klinreq.ReqInfo{
		Dest:       host,
		Dport:      klinutils.Stringtoport("govirthost"),
		Method:     "GET",
		Payload:    p,
		TrustBytes: c.Tb,
		CertBytes:  c.Cb,
		KeyBytes:   c.Kb,
	}
	resp, err := klinreq.SendPayload(i)
	if err != nil {
		return r, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(string(body))
		return r, errors.New("Failed, check logs on the govirthost server" + host)
	}
	var tmpr govirtlib.ReturnPayload
	err = json.Unmarshal(body, &tmpr)
	if err != nil {
		return r, err
	}
	return tmpr.Xml, err
}
func (c *Conn) getvms(hosts []string) (govirtlib.ReturnPayload, error) {
	var r govirtlib.ReturnPayload
    listvms := make(map[string][]govirtlib.DomainInfo)
	p := &govirtlib.GetPayload{
		Target: "vm",
	}
	for _, host := range hosts {
		i := &klinreq.ReqInfo{
			Dest:       host,
			Dport:      klinutils.Stringtoport("govirthost"),
			Method:     "GET",
			Payload:    p,
			TrustBytes: c.Tb,
			CertBytes:  c.Cb,
			KeyBytes:   c.Kb,
		}
		resp, err := klinreq.SendPayload(i)
		if err != nil {
			return r, err
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return r, err
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			fmt.Println(string(body))
			return r, errors.New("Failed, check logs on the govirthost server" + host)
		}
		var tmpr govirtlib.ReturnPayload
		err = json.Unmarshal(body, &tmpr)
		if err != nil {
			return r, err
		}
        listvms[host] = tmpr.Domains
	}
    r.Listvms = listvms
	return r, nil
}
