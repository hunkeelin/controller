package controller

import (
	"fmt"
    "io/ioutil"
    "testing"
    "github.com/hunkeelin/SuperCAclient/lib"
    "github.com/hunkeelin/pki"
    "net/http"
    "github.com/hunkeelin/klinutils"
    "github.com/hunkeelin/mtls/klinserver"
)
func TestServer(t *testing.T) {
    fmt.Println("testing Server")
    c := Conn {
    }
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
        Route: "cacerts/" + "govirt"+ ".crt",
    }
    c.Tb, err = klinutils.Wget(r)
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
    c.Cb, c.Kb, err = client.Getkeycrtbyte(w)
    if err != nil {
        panic(err)
    }
    m, err := Parse("config")
    if err != nil {
        panic(err)
    }
    c.Clusters = m
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
    j := &klinserver.ServerConfig {
        BindPort: klinutils.Stringtoport("controller"),
        ServeMux: con,
        Https: false,
    }
    insecure := false
    if !insecure {
        j.CertBytes = c.Cb
        j.KeyBytes = c.Kb
        j.Https = true
        j.Verify = false
        j.TrustBytes = c.Tb
    }
    err = klinserver.Server(j)
    if err != nil {
        panic(err)
    }
}
