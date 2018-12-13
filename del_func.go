package controller

import (
    "encoding/hex"
    "fmt"
    "encoding/base64"
    "crypto/sha256"
    "strings"
    "errors"
	"github.com/hunkeelin/govirt/govirtlib"
	"net/http"
	"strconv"
)

func (c *Conn) delvm(w http.ResponseWriter, r *http.Request, v govirtlib.PostPayload) error {
    var err error
    // vm to delete
    d, _ := base64.StdEncoding.DecodeString(r.Header.Get("api-key"))
    userpw := strings.Split(string(d), ":")
    usersum := sha256.Sum256([]byte(userpw[0]))
    userhash := hex.EncodeToString(usersum[:])
    todelete := v.Domain
    for _,vhosts := range c.Clusters {
        p, err := c.getvms(vhosts.Govirt)
        if err != nil {
            return err
        }
        for parent, hosts := range p.Listvms {
            for _,i := range hosts {
                if i.Domain.Name == todelete {
                    vmhash := hex.EncodeToString(i.Domain.UUID[:])
                    if vmhash[0:8] != userhash[0:8] {
                        fmt.Println("The vm",i.Domain.Name,"doesn't belong to",userpw[0])
                        return errors.New("Unable to delete " +i.Domain.Name)
                    }
                    tmp,_ := c.userlimit[userhash[0:8]]
                    ocpu, err := strconv.Atoi(vmhash[8:10])
                    if err != nil {
                        return err
                    }
                    oram, err := strconv.Atoi(vmhash[10:12])
                    if err != nil {
                        return err
                    }
                    if i.State == "running" {
                        err = c.statevm("destroy",todelete,parent)
                        if err != nil {
                            return err
                        }
                    }
                    err = c.statevm("undefine",todelete,parent)
                    if err != nil {
                        return err
                    }
                    tmp.vcpu += ocpu
                    tmp.vram += oram
                    c.userlimit[userhash[0:8]] = tmp
                }
            }
        }
    }
    err = c.delhost_network(c.Clusters[v.Cluster].Godhcp, todelete)
    if err != nil {
        return err
    }
    err = c.delimage(c.Clusters[v.Cluster].Storage,todelete)
    if err != nil {
        return err
    }
	return nil
}
