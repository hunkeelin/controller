package controller

import (
	"sync"
	"time"
)
type ClusterInfo struct {
    ClusterName string
    Godhcp      string
    Govirt      []string
    Storage     string
}
type Conn struct {
	Cb     []byte
	Kb     []byte
	Tb     []byte
	postMu sync.Mutex
	authcb []byte
	authkb []byte
	authtb []byte
	Ixml   map[string][]byte
	rmap   map[string]rlimit
    Clusters map[string]ClusterInfo
}
type rlimit struct {
	cpu       int       // vcpu
	mem       int       // mem in GB
	timelimit time.Time // in hours.h
}
