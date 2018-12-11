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
