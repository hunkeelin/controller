package controller

import (
	"errors"
	"fmt"
	"github.com/hunkeelin/govirt/govirtlib"
	"io/ioutil"
	"net/http"
	"strings"
)

func (c *Conn) post(w http.ResponseWriter, r *http.Request) error {
	var p govirtlib.PostPayload
    var err error
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("unable to read response body post")
		return err
	}
	err = json.Unmarshal(b, &p)
	if err != nil {
		fmt.Println("unable to unmarshal json post")
		return err
	}
    if _,ok := c.Clusters[p.Cluster]; !ok {
        return errors.New("Cluster "+p.Cluster+" doens't exist\n")
    }
	switch strings.ToLower(p.Action) {
	case "createvm":
		err = c.createvm(w,r,p)
	default:
		return errors.New("Invalid Action")
	}
	return err
}
