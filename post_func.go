package controller

import (
    "fmt"
	"bytes"
    "encoding/hex"
    "crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/hunkeelin/govirt/govirtlib"
	"github.com/hunkeelin/klinutils"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (c *Conn) Migratehost(ori, dest, vm string) error {
	err := c.migrate(ori, dest, vm)
	if err != nil {
		return err
	}
	xml, err := c.getxml(vm, ori)
	if err != nil {
		return err
	}
	err = c.statevm("unDefine", vm, ori)
	if err != nil {
		return err
	}
	err = c.Define(xml, dest)
	if err != nil {
		return err
	}
	return nil
}
func checkVmForm(v govirtlib.CreateVmForm) error {
	switch {
	case v.Hostname == "":
		return errors.New("Please specify hostname")
	case v.Uuid == "":
		return errors.New("Please specifiy uuid")
	case v.MemoryCount == 0:
		return errors.New("Please specify memorycount")
	case v.CpuCount == 0:
		return errors.New("Please specify memorycount")
	case !klinutils.Is_mac(v.VmMac):
		return errors.New("Please specify a valid mac address")
	case v.Vlan == "":
		return errors.New("Please specify Vlan for network")
	case !klinutils.Is_ipv4(v.VmIp):
		return errors.New("Please speceify a valid IP")
	default:
		return nil
	}
	return nil
}
func (c *Conn) createvm(w http.ResponseWriter, r *http.Request, v govirtlib.PostPayload) error {
	d, _ := base64.StdEncoding.DecodeString(r.Header.Get("api-key"))
	userpw := strings.Split(string(d), ":")
    usersum := sha256.Sum256([]byte(userpw[0]))
    userhash := hex.EncodeToString(usersum[:])
    resource := c.userlimit[userhash[0:8]]
	if resource.vcpu < v.VmForm.CpuCount {
		rcp := strconv.Itoa(resource.vcpu)
		return errors.New(userpw[0] + " Exceed cpu quota you have " + rcp)
	}
	if resource.vram < v.VmForm.MemoryCount {
		rrm := strconv.Itoa(resource.vram)
		return errors.New(userpw[0] + " Exceed mem quota you have " + rrm)
	}
	err := checkVmForm(v.VmForm)
	if err != nil {
		return err
	}
	if c.Ixml[v.VmForm.Image] == nil {
		return errors.New("No image for : " + v.VmForm.Image)
	}
	err = c.setimage(c.Clusters[v.Cluster].Storage, v.VmForm.Image, v.VmForm.Hostname)
	if err != nil {
		return err
	}
	err = c.edithost(c.Clusters[v.Cluster].Godhcp, v, false)
	if err != nil {
		return err
	}
	uuid, err := klinutils.Genuuidv2(userpw[0], v.VmForm.CpuCount, v.VmForm.MemoryCount)
	if err != nil {
		return err
	}
	xml := c.Ixml[v.VmForm.Image]
	xml = bytes.Replace(xml, []byte("name_replace"), []byte(v.VmForm.Hostname), -1)
	xml = bytes.Replace(xml, []byte("uuid_replace"), uuid, -1)
	xml = bytes.Replace(xml, []byte("memory_replace"), []byte(strconv.Itoa(v.VmForm.MemoryCount)), -1)
	xml = bytes.Replace(xml, []byte("cpu_replace"), []byte(strconv.Itoa(v.VmForm.CpuCount)), -1)
	xml = bytes.Replace(xml, []byte("imagedir_replace"), []byte("/data/govirt/storage"), -1)
	xml = bytes.Replace(xml, []byte("mac_replace"), []byte(v.VmForm.VmMac), -1)
	xml = bytes.Replace(xml, []byte("vlan_replace"), []byte(v.VmForm.Vlan), -1)
	rand.Seed(time.Now().UTC().UnixNano())
	randhostint := klinutils.RandInt(0, len(c.Clusters[v.Cluster].Govirt))
	err = c.Define(xml, c.Clusters[v.Cluster].Govirt[randhostint])
	if err != nil {
		panic(err)
	}
	err = c.statevm("start", v.VmForm.Hostname, c.Clusters[v.Cluster].Govirt[randhostint])
	if err != nil {
		panic(err)
	}
	resource.vcpu = resource.vcpu - v.VmForm.CpuCount
	resource.vram = resource.vram - v.VmForm.MemoryCount
	c.userlimit[userhash[0:8]] = resource
	return nil
}
func (c *Conn) CreateNewVm(v govirtlib.PostPayload) error {
	err := checkVmForm(v.VmForm)
	if err != nil {
		return err
	}
	if c.Ixml[v.VmForm.Image] == nil {
		return errors.New("No image for : " + v.VmForm.Image)
	}
	m, err := Parse("config")
	if err != nil {
		panic(err)
	}
	err = c.edithost(m[v.Cluster].Godhcp, v, false)
	if err != nil {
		panic(err)
	}
	err = c.setimage(m[v.Cluster].Storage, v.VmForm.Image, v.VmForm.Hostname)
	if err != nil {
		panic(err)
	}
	xml := c.Ixml[v.VmForm.Image]
	xml = bytes.Replace(xml, []byte("name_replace"), []byte(v.VmForm.Hostname), -1)
	xml = bytes.Replace(xml, []byte("uuid_replace"), []byte(v.VmForm.Uuid), -1)
	xml = bytes.Replace(xml, []byte("memory_replace"), []byte(strconv.Itoa(v.VmForm.MemoryCount)), -1)
	xml = bytes.Replace(xml, []byte("cpu_replace"), []byte(strconv.Itoa(v.VmForm.CpuCount)), -1)
	xml = bytes.Replace(xml, []byte("imagedir_replace"), []byte("/data/govirt/storage"), -1)
	xml = bytes.Replace(xml, []byte("mac_replace"), []byte(v.VmForm.VmMac), -1)
	xml = bytes.Replace(xml, []byte("vlan_replace"), []byte(v.VmForm.Vlan), -1)
	rand.Seed(time.Now().UTC().UnixNano())
	randhostint := klinutils.RandInt(0, len(m[v.Cluster].Govirt))
	err = c.Define(xml, m[v.Cluster].Govirt[randhostint])
	if err != nil {
		panic(err)
	}
	err = c.statevm("start", v.VmForm.Hostname, m[v.Cluster].Govirt[randhostint])
	if err != nil {
		panic(err)
	}
	return nil
}

func (c *Conn) statevmapi(w http.ResponseWriter, r *http.Request, v govirtlib.PostPayload) error {
    allowops := []string{"shutdown","start","destroy"}
    if !klinutils.StringInSlice(v.Target,allowops){
        return errors.New("Please specify correct target operations shutdown/start/destory")
    }
    d, _ := base64.StdEncoding.DecodeString(r.Header.Get("api-key"))
    userpw := strings.Split(string(d), ":")
    usersum := sha256.Sum256([]byte(userpw[0]))
    userhash := hex.EncodeToString(usersum[:])
    tostate := v.Domain
    for _,vhosts := range c.Clusters {
        p, err := c.getvms(vhosts.Govirt)
        if err != nil {
            return err
        }
        for parent, hosts := range p.Listvms {
            for _,i := range hosts {
                if i.Domain.Name == tostate {
                    vmhash := hex.EncodeToString(i.Domain.UUID[:])
                    if vmhash[0:8] != userhash[0:8] {
                        fmt.Println("The vm",i.Domain.Name,"doesn't belong to",userpw[0])
                        return errors.New("Unable to "+v.Target+" "+i.Domain.Name)
                    }
                    err = c.statevm(v.Target,tostate,parent)
                    if err != nil {
                        return err
                    }
                }
            }
        }
    }
    return nil
}
