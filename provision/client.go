package provision
/*
import  (
  "fmt"
  "http"
)

func LogsRequest(containerid string, containername string) error {
	gulpUrl := "http://localhost:6666/docker/logs"
	url := gulpUrl + "docker/logs"
	fmt.Println("URL:>", url)

	data := &global.DockerLogsInfo{ContainerId: containerid, ContainerName: containername}
	res2B, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(res2B))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return nil
}




func NetworkRequest(containerid string, ip string) error {
	url := "http://localhost:6666/docker/networks"
    fmt.Println("URL:>", url)
    data := &global.DockerNetworksInfo{Bridge: "one", ContainerId: containerid, IpAddr: ip, Gateway: "103.56.92.1"}
	res2B, _ := json.Marshal(data)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(res2B))
    req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
    return nil
}
*/
