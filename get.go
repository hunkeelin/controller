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
	switch strings.ToLower(p.Target) {
//	case "vms":
//		err = c.apigetvms(w, r)
//		if err != nil {
//			fmt.Println("Unable to get images")
//			return err
//		}
	default:
		return errors.New("Invalid Storage get Action " + strings.ToLower(p.Target))
	}
	return nil
}
