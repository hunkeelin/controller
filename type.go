package controller

import (
	"sync"
	"time"
)
type authData struct {
    User string `json:"user"`
    Password string `json:"password"`
    ValidGroups []string `json:"validgroups"`
    ValidUsers []string `json:"validusers"`
}
type ClusterInfo struct {
    ClusterName string
    Godhcp      string
    Govirt      []string
    Storage     string
}
type Conn struct {
	cb     []byte
	kb     []byte
	tb     []byte
	postMu sync.Mutex
	authcb []byte
	authkb []byte
	authtb []byte
	Ixml   map[string][]byte
	rmap   map[string]rlimit
    Clusters map[string]ClusterInfo
    userlimit map[string]resourcelimit
}
type resourcelimit struct {
    vcpu int
    vram int
    active bool
}
type rlimit struct {
	cpu       int       // vcpu
	mem       int       // mem in GB
	timelimit time.Time // in hours.h
}
type Payload struct {
    Rrsets `json:"rrsets"`
}
type Rrsets []struct {
    Name       string `json:"name"`
    Type       string `json:"type"`
    TTL        int    `json:"ttl"`
    Changetype string `json:"changetype"`
    Records `json:"records"`
}
type Records []struct {
    Content  string `json:"content"`
    Disabled bool `json:"disabled"`
}
func dnsPayload(name,ptype,ctype,content string ,ttl int, disable bool) Payload{
    records := Records{
        {content,disable},
    }
    rrsets := Rrsets{
        {name,ptype,ttl,ctype,records},
    }
    return Payload{rrsets}
}
