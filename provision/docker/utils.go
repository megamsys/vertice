/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*func IPRequest(subnet net.IPNet) (net.IP, uint, error) {
	bits := bitCount(subnet)
	bc := int(bits / 8)
	partial := int(math.Mod(bits, float64(8)))
	if partial != 0 {
		bc += 1
	}
	index := global.IPIndex{}
	res, err := index.Get(global.IPINDEXKEY)
	if err != nil {
		log.Errorf("Failed to load %s from riak: %s.", global.IPINDEXKEY, err.Error())
		return nil, 0, err
	}

	return getIP(subnet, res.Index+1), res.Index + 1, nil
}

// Given Subnet of interest and free bit position, this method returns the corresponding ip address
// This method is functional and tested. Refer to ipam_test.go But can be improved

func getIP(subnet net.IPNet, pos uint) net.IP {
	retAddr := make([]byte, len(subnet.IP))
	copy(retAddr, subnet.IP)

	mask, _ := subnet.Mask.Size()
	var tb, byteCount, bitCount int
	if subnet.IP.To4() != nil {
		tb = 4
		byteCount = (32 - mask) / 8
		bitCount = (32 - mask) % 8
	} else {
		tb = 16
		byteCount = (128 - mask) / 8
		bitCount = (128 - mask) % 8
	}

	for i := 0; i <= byteCount; i++ {
		maskLen := 0xFF
		if i == byteCount {
			if bitCount != 0 {
				maskLen = int(math.Pow(2, float64(bitCount))) - 1
			} else {
				maskLen = 0
			}
		}
		masked := pos & uint((0xFF&maskLen)<<uint(8*i))
		retAddr[tb-i-1] |= byte(masked >> uint(8*i))
	}
	return net.IP(retAddr)
}

func bitCount(addr net.IPNet) float64 {
	mask, _ := addr.Mask.Size()
	if addr.IP.To4() != nil {
		return math.Pow(2, float64(32-mask))
	} else {
		return math.Pow(2, float64(128-mask))
	}
}

func testAndSetBit(a []byte) uint {
	var i uint
	for i = uint(0); i < uint(len(a)*8); i++ {
		if !testBit(a, i) {
			setBit(a, i)
			return i + 1
		}
	}
	return i
}

func testBit(a []byte, k uint) bool {
	return ((a[k/8] & (1 << (k % 8))) != 0)
}

func setBit(a []byte, k uint) {
	a[k/8] |= 1 << (k % 8)
}

func setContainerNAL(container *global.Container) (string, error) {
	//generate the ip
	subnetip, _ := config.GetString("docker:subnet")
	_, subnet, _ := net.ParseCIDR(subnetip)
	ip, pos, err := IPRequest(*subnet)
	if err != nil {
		log.Errorf("Failed to generate an ip address from docker subnet: %s", err.Error())
		return "", err
	}
	client, _ := docker.NewClient("http://" + container.SwarmNode + ":2375")
	ch := make(chan bool)

//configure ip to container
	go recv(container, ip.String(), client, ch)

	err := updateIndex(ip.String(), pos)
	if err != nil {
		log.Errorf("Failed to increment the ip in riak: %s", err.Error())
	}
	return ip.String(), nil
}



// UpdateComponent updates the ipaddress that is bound to the container
//It talks to riakdb and updates the respective component(s)
func updateContainerJSON(assembly *app.DeepAssembly, container *global.Container, endpoint string) {

	log.Debugf("[docker] updating container")
	mySlice := make([]*global.KeyValuePair, 5)
	mySlice[0] = &global.KeyValuePair{Key: "ip", Value: container.IPAddress}
	mySlice[1] = &global.KeyValuePair{Key: "id", Value: container.ContainerID}
	mySlice[2] = &global.KeyValuePair{Key: "port", Value: container.Ports}
	mySlice[3] = &global.KeyValuePair{Key: "endpoint", Value: endpoint}
	mySlice[4] = &global.KeyValuePair{Key: "host", Value: container.SwarmNode}

	update := global.Component{
		Id:                assembly.Components[0].Id,
		Name:              assembly.Components[0].Name,
		ToscaType:         assembly.Components[0].ToscaType,
		Inputs:            assembly.Components[0].Inputs,
		Outputs:           mySlice,
		Artifacts:         assembly.Components[0].Artifacts,
		RelatedComponents: assembly.Components[0].RelatedComponents,
		Operations:        assembly.Components[0].Operations,
		Status:            assembly.Components[0].Status,
		CreatedAt:         assembly.Components[0].CreatedAt,
	}

	conn, err := db.Conn("components")
	if err != nil {
		log.Errorf("Failed to connect to riak: %s", err.Error())
	}

	err := conn.StoreStruct(assembly.Components[0].Id, &update)
	if err != nil {
		log.Errorf("Failed to store the component in riak: %s", err.Error())
	}
	log.Debugf("[docker] Componet updated successfully.")
}

func GetMemory() int64 {
	memory, _ := config.GetInt("docker:memory")
	return int64(memory)
}

func GetSwap() int64 {
	swap, _ := config.GetInt("docker:swap")
	return int64(swap)

}

func GetCpuPeriod() int64 {
	cpuPeriod, _ := config.GetInt("docker:cpuperiod")
	return int64(cpuPeriod)

}

func GetCpuQuota() int64 {
	cpuQuota, _ := config.GetInt("docker:cpuquota")
	return int64(cpuQuota)

}

//this shall be configurabel parameter saying wait_container_up = 18000 (ms)
func recv(container *global.Container, ip string, client *docker.Client, ch chan bool) {
	log.Debugf("[docker] wait for the container to start 18000 ms")
	time.Sleep(18000 * time.Millisecond)


	 //Inspect API is called to fetch the data about the launched container
	inscontainer, _ := client.InspectContainer(container.ContainerID)
	contain := &docker.Container{}
	mapC, _ := json.Marshal(inscontainer)
	json.Unmarshal([]byte(string(mapC)), contain)

	container_state := &docker.State{}
	mapN, _ := json.Marshal(contain.State)
	json.Unmarshal([]byte(string(mapN)), container_state)

	if container_state.Running == true {
		postnetwork(container, ip)
		postlogs(container)
		ch <- true
		return
	}

	go recv(container, ip, client, ch)
}

func postnetwork(container *global.Container, ip string) {
	gulpPort, _ := config.GetInt("docker:gulp_port")

	url := "http://" + container.SwarmNode + ":" + strconv.Itoa(gulpPort) + "/docker/networks"
	log.Debugf("[docker] URL:> %s", url)

	bridge, _ := config.GetString("docker:bridge")
	gateway, _ := config.GetString("docker:gateway")

	data := &global.DockerNetworksInfo{Bridge: bridge, ContainerId: container.ContainerID, IpAddr: ip, Gateway: gateway}
	res2B, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(res2B))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Failed to connect to megamdocker node client: %s", err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf(`[docker]
		*------------------------------------------------------*
	  HTTP   :   %s      STATUS : %s
		HEADER :   %s
		BODY   :   %s
  	*-------------------------------------------------------*`,
		url, resp.Status, resp.Header, string(body))
}

func postlogs(container *global.Container) error {
	gulpPort, _ := config.GetInt("docker:gulp_port")
	url := "http://" + container.SwarmNode + ":" + strconv.Itoa(gulpPort) + "/docker/logs"
	log.Debugf("[docker] URL:> %s", url)

	data := &global.DockerLogsInfo{ContainerId: container.ContainerID, ContainerName: container.ContainerName}
	res2B, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(res2B))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Failed to connect to megamdocker node client: %s", err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf(`*------------------------------------------------------*
	            HTTP   :   %s      STATUS : %s
							HEADER :   %s
							BODY   :   %s
							*-------------------------------------------------------*`,
							url, resp.Status, resp.Header, string(body))
	return nil
}

func updateIndex(ip string, pos uint) error {
	index := global.IPIndex{}
	res, err := index.Get(global.IPINDEXKEY)
	if err != nil {
		log.Errorf("Failed to store in riak: %s.", err.Error())
		return err
	}

	update := global.IPIndex{
		Ip:     ip,
		Subnet: res.Subnet,
		Index:  pos,
	}

	conn, err := db.Conn("ipindex")
	if err != nil {
		log.Errorf("Failed to connect to riak: %s", err.Error())
		return err
	}

	err := conn.StoreStruct(global.IPINDEXKEY, &update)
	if err != nil {
		log.Error("Failed to store the updated index value : %s", err.Error())
		return serr
	}
	log.Debugf("[docker] Docker network index updated successfully.")
	return nil
}


*/
