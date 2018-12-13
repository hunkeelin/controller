package controller

import (
	"errors"
    "net/http"
    "strings"
	"fmt"
    "encoding/hex"
    "encoding/base64"
    "crypto/sha256"
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
		TrustBytes: c.tb,
		CertBytes:  c.cb,
		KeyBytes:   c.kb,
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
			TrustBytes: c.tb,
			CertBytes:  c.cb,
			KeyBytes:   c.kb,
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
func (c *Conn) getnetapi(w http.ResponseWriter, nethost string) error {
    payload := &govirtlib.GetPayload{
        Target: "network",
    }
    i := &klinreq.ReqInfo{
        Dest:    nethost,
        Dport:   klinutils.Stringtoport("godhcp"),
        Method:  "GET",
        Payload: payload,
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
        return errors.New("Failed, check logs on the godhcp server" + nethost)
    }
    var p govirtlib.ReturnPayload
    err = json.Unmarshal(body, &p)
    if err != nil {
        return err
    }
    err = json.NewEncoder(w).Encode(p)
    if err != nil {
        fmt.Println("unable to encode json")
        return err
    }
    return nil
}
func (c *Conn) getvmsapi(w http.ResponseWriter,r *http.Request,hosts []string) error {
	var rp govirtlib.ReturnPayload
    d, _ := base64.StdEncoding.DecodeString(r.Header.Get("api-key"))
    userpw := strings.Split(string(d), ":")
    usersum := sha256.Sum256([]byte(userpw[0]))
    userhash := hex.EncodeToString(usersum[:])
    var udomainlist []govirtlib.DomainInfo
	p := &govirtlib.GetPayload{
		Target: "vm",
	}
	for _, host := range hosts {
		i := &klinreq.ReqInfo{
			Dest:       host,
			Dport:      klinutils.Stringtoport("govirthost"),
			Method:     "GET",
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
			return errors.New("Failed, check logs on the govirthost server" + host)
		}
		var tmpr govirtlib.ReturnPayload
		err = json.Unmarshal(body, &tmpr)
		if err != nil {
			return err
		}
        for _,val := range tmpr.Domains{
            vmhash := hex.EncodeToString(val.Domain.UUID[:])
            if vmhash[0:8] == userhash[0:8] {
                udomainlist = append(udomainlist,val)
            }
        }
	}
    rp.Domains = udomainlist
    err := json.NewEncoder(w).Encode(rp)
    if err != nil {
        fmt.Println("unable to encode json")
        return err
    }
	return nil
}
