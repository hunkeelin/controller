# Introduction #
The controller for govirt

## Authentication 
BASIC authentication of header `api-key` and base64 string of "$ldapuser:$ldappassword" including the colon.

## GET 
``` 
{
    "Target": $target
}
```

### target = vms ###
Will list all the vms belonging to the user. You will get a json payload in terms of ReturnPayload. 

```
type ReturnPayload struct {
    Domains        []DomainInfo            `json:"domains"`
}
type DomainInfo struct {
    Domain libvirt.Domain `json:"domain"`
    State  string         `json:"state"`
}
type libvirt.Domain {
    Name string
    UUID UUID
    ID   int32 
}
```


## POST 
```
type PostPayload struct {
    Target        string         `json:"target"`
    Cluster       string         `json:"cluster"`
    Domain        string         `json:"domain"`
    Action        string         `json:"action"`
    VmForm        CreateVmForm   `json:"createvmform"`
}
```

### Action = createvm ###
This action for create VM. You will have to include an additional payload `VmForm` shown below. 

Require frields: 
    - Cluster: Specify the govirt cluster. E.g `sf_deploy`
    - VmForm:
```
type CreateVmForm struct {
    Hostname    string `json:"hostname"`
    VmIp        string `json:"ip"`          // the ip address of the virutal machine
    CpuCount    int    `json:"cpucount"`    // cpucount
    MemoryCount int    `json:"memorycount"` // in GB
    Image       string `json:"image"`       // the image to use
}
```
### Action = state ###
Manupulate the state of the vm.

Require fields: 
    - Cluster: Specify the govirt cluster. E.g `sf_deploy`
    - Target: The desire state of the virtual machine. start/destroy/shutdown
    - Domain: The virtual machine name. 

## DELETE 
```
type PostPayload struct {
    Domain        string         `json:"action"
    Target        string         `json:"target"`
}
```
### Target = vm ###
To remove the virtual machine from the cluster 

Require fields: 
     - Domain: The name of the virutal machine you want to delete. 

##CLI
For cli please read `_test.go`
