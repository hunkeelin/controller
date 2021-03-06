package controller

import (
	"fmt"
    "encoding/hex"
	"github.com/hunkeelin/SuperCAclient/lib"
	"github.com/hunkeelin/klinutils"
    "strconv"
	"github.com/hunkeelin/mtls/klinserver"
	"github.com/hunkeelin/pki"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	fmt.Println("testing Server")
	c := Conn{}
	r := klinutils.WgetInfo{
		Dest:  "ec2-superca-prod-1.squaretrade.com",
		Dport: "2018",
		Route: "cacerts/rootca.crt",
	}
	cab, err := klinutils.Wget(r)
	if err != nil {
		panic(err)
	}
	r = klinutils.WgetInfo{
		Dest:  "ec2-superca-prod-1.squaretrade.com",
		Dport: "2018",
		Route: "cacerts/" + "govirt" + ".crt",
	}
	c.tb, err = klinutils.Wget(r)
	if err != nil {
		panic(err)
	}
	r.Route = "cacerts/superauth.crt"
	c.authtb, err = klinutils.Wget(r)
	if err != nil {
		panic(err)
	}
	w := client.WriteInfo{
		CAName:  "ec2-superca-prod-1.squaretrade.com",
		CABytes: cab,
		CAport:  "2018",
		Chain:   true,
		CSRConfig: &klinpki.CSRConfig{
			RsaBits: 2048,
		},
		SignCA: "govirtcon",
	}
	c.cb, c.kb, err = client.Getkeycrtbyte(w)
	if err != nil {
		panic(err)
	}
	w.SignCA = "superauth"
	c.authcb, c.authkb, err = client.Getkeycrtbyte(w)
	if err != nil {
		panic(err)
	}
	m, err := Parse("config")
	if err != nil {
		panic(err)
	}
	c.Clusters = m
	rlim := make(map[string]resourcelimit)
    for _, vhost := range m {
        pp, err := c.getvms(vhost.Govirt)
        if err != nil {
            panic(err)
        }
        for _, hosts := range pp.Listvms {
            for _, i := range hosts {
                userhash := hex.EncodeToString(i.Domain.UUID[:])
                tmp,_ := rlim[userhash[0:8]]
                ocpu, err := strconv.Atoi(userhash[8:10])
                if err != nil {
                    fmt.Println("this host does not belong to a user",i.Domain.Name)
                    continue
                }
                oram, err := strconv.Atoi(userhash[10:12])
                if err != nil {
                    fmt.Println("this host does not belong to a user",i.Domain.Name)
                    continue
                }
                tmp.vcpu = tmp.vcpu - ocpu
                tmp.vram = tmp.vram - oram
                rlim[userhash[0:8]] = tmp
            }
        }
    }
	c.userlimit = rlim
	con := http.NewServeMux()
	con.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c.MainHandler(w, r)
	})
	ctemp, err := ioutil.ReadFile("ctemplate.xml")
	if err != nil {
		panic(err)
	}
	utemp, err := ioutil.ReadFile("ctemplate.xml")
	if err != nil {
		panic(err)
	}
	ixml := make(map[string][]byte)
	ixml["ubuntu"] = utemp
	ixml["centos"] = ctemp
	c.Ixml = ixml
	j := &klinserver.ServerConfig{
		BindPort: klinutils.Stringtoport("controller"),
		ServeMux: con,
		Https:    false,
	}
	insecure := false
	if !insecure {
		j.CertBytes = c.cb
		j.KeyBytes = c.kb
		j.Https = true
		j.Verify = false
		j.TrustBytes = c.tb
	}
	err = klinserver.Server(j)
	if err != nil {
		panic(err)
	}
}
