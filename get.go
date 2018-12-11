package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
    "github.com/hunkeelin/govirt/govirtlib"
	"net/http"
	"strings"
)

func (c *Conn) get(w http.ResponseWriter, r *http.Request) error {
	var p govirtlib.GetPayload
    var err error
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("unable to read response body post")
	}
	err = json.Unmarshal(b, &p)
    if _,ok := c.Clusters[p.Cluster]; !ok {
        return errors.New("Cluster "+p.Cluster+" doens't exist\n")
    }
	switch strings.ToLower(p.Target) {
	case "vms":
		err = c.getvmsapi(w, c.Clusters[p.Cluster].Govirt)
	case "network":
		err = c.getnetapi(w, c.Clusters[p.Cluster].Godhcp)
	default:
		return errors.New("Invalid Storage get Action " + strings.ToLower(p.Target))
	}
	return err
}
