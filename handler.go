package controller

import (
	"fmt"
    "io/ioutil"
    "strings"
    "github.com/hunkeelin/mtls/klinreq"
    "encoding/base64"
	"github.com/json-iterator/go"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (c *Conn) MainHandler(w http.ResponseWriter, r *http.Request) {
    var err error
    d,err := base64.StdEncoding.DecodeString(r.Header.Get("api-key"))
    if err != nil {
        fmt.Println("Unable to decode given api-key is it base64 encoded?")
        w.WriteHeader(500)
        w.Write([]byte(err.Error()))
        return
    }
    userpw := strings.Split(string(d),":")
    if len(userpw) != 2 {
        fmt.Println("did user provide a encoded64 string?")
        w.WriteHeader(401)
        w.Write([]byte("please provide api-key\n"))
        return
    }
    a,b := c.checkldap(userpw[0],userpw[1],[]string{"it","engineering"})
    if b != 200 {
        fmt.Println(userpw[0],"user not authorized",a)
        w.WriteHeader(401)
        w.Write([]byte("user not authorized"))
        return
    }
    resource := c.userlimit[userpw[0]]
    if !resource.active {
        resource.vcpu = 8
        resource.vram = 16
        resource.active = true
    }
    c.userlimit[userpw[0]] = resource
	switch r.Method {
	case "GET":
		err = c.get(w, r)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
		}
	case "POST":
		err = c.post(w, r)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
		}
//	case "DELETE":
//		err := c.del(w, r)
//		if err != nil {
//			fmt.Println(err)
//			w.WriteHeader(500)
//			w.Write([]byte(err.Error()))
//		}
	default:
		fmt.Println("invalid method")
		w.WriteHeader(500)
	}
    if err != nil {
        w.Write([]byte(err.Error()))
    }
    return
}

func (c *Conn) checkldap (user,pw string,v []string) (string,int){
    p := authData {
        User: user,
        Password: pw,
        ValidGroups: v,
    }
    i := &klinreq.ReqInfo {
        Dest: "ec2-auth-prod-1.squaretrade.com",
        Dport: "2014",
        Route: "ldap",
        Payload: p,
        HttpVersion: 1,
        TimeOut: 1500,
        CertBytes: c.authcb,
        KeyBytes: c.authkb,
        TrustBytes: c.authtb,
    }
    resp, err := klinreq.SendPayload(i)
    if err != nil {
        fmt.Println(err)
        return "Server Error",500
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Unable to read response body from ldapapi")
        return "Server Error",500
    }
    resp.Body.Close()
    return string(body),resp.StatusCode
}
