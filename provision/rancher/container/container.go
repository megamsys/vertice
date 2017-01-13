package container

import (
	//"fmt"
	//"io"
	"net"
	"net/url"
	"time"
	//"bytes"
	//	"os"
	//	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/go-rancher/v2"
	"github.com/megamsys/libgo/utils"
	//"github.com/megamsys/libgo/events"
	//	"github.com/megamsys/libgo/events/alerts"
	"github.com/megamsys/libgo/safe"
	"github.com/megamsys/vertice/carton"
	"github.com/megamsys/vertice/provision"
	"github.com/megamsys/vertice/provision/rancher/cluster"
)

const (
	portRangeStart    = 49153
	portRangeEnd      = 65535
	portAllocMaxTries = 15
)

type RancherProvisioner interface {
	Cluster() *cluster.Cluster
	//	PushImage(name, tag string) error
}

type Container struct {
	Id                      string //container id.
	BoxId                   string
	AccountId               string
	CartonId                string
	Name                    string
	BoxName                 string
	Level                   provision.BoxLevel
	PublicIp                string
	HostId                  string
	HostAddr                string
	HostPort                string
	PrivateKey              string
	Client                  *client.RancherClient
	Version                 string
	Image                   string
	Status                  utils.Status
	State                   utils.State
	BuildingImage           string
	LastStatusUpdate        time.Time
	LastSuccessStatusUpdate time.Time
	LockedUntil             time.Time
	Routable                bool
	Region                  string
	closechan               chan bool
}

/*
func (c *Container) ShortId() string {
	if len(c.Id) > 10 {
		return c.Id[:10]
	}
	return c.Id
}

func (c *Container) Available() bool {
	return c.Status.String() == constants.StatusContainerStarted.String() || c.Status.String() == constants.StatusContainerStarting.String()
}

func (c *Container) Address() *url.URL {
	return &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", c.HostAddr, c.HostPort),
	}
}
*/
type CreateArgs struct {
	ImageId          string
	Box              *provision.Box
	Deploy           bool
	Provisioner      RancherProvisioner
	DestinationHosts []string
}

func (c *Container) Create(args *CreateArgs) error {
	config := client.Container{
		Name:          c.BoxName,
		ImageUuid:     args.Box.Repo.Source + ":" + args.ImageId,
		Memory:        int64(args.Box.ConGetMemory()),
		MemorySwap:    int64(args.Box.ConGetMemory() + args.Box.GetSwap()),
		CpuShares:     int64(args.Box.GetCpushare()),
		StartOnCreate: true,
	}

	cl := args.Provisioner.Cluster()
	cl.Region = args.Box.Region

	addr, cont, err := cl.CreateContainerSchedulerOpts(config)
	if err != nil {
		log.Errorf("Error on creating container in docker %s - %s", c.BoxName, err)
		return err
	}
	c.Id = cont.Id
	c.HostAddr = urlToHost(addr)
	return nil
}
func (c *Container) SetStatus(status utils.Status) error {
	log.Debugf("  set status[%s] of container (%s, %s)", c.BoxId, c.Name, status.String())
	if asm, err := carton.NewAssembly(c.CartonId, c.AccountId, ""); err != nil {
		return err
	} else if err = asm.SetStatus(status); err != nil {
		return err
	}

	if c.Level == provision.BoxSome {
		log.Debugf("  set status[%s] of container (%s, %s)", c.BoxId, c.Name, status.String())
		if comp, err := carton.NewComponent(c.BoxId, c.AccountId, ""); err != nil {
			return err
		} else if err = comp.SetStatus(status, c.AccountId, ""); err != nil {
			return err
		}
	}
	return nil
}

func (c *Container) SetMileStone(state utils.State) error {
	log.Debugf("  set state[%s] of container (%s, %s)", c.BoxId, c.Name, state.String())
	if asm, err := carton.NewAssembly(c.CartonId, c.AccountId, ""); err != nil {
		return err
	} else if err = asm.SetState(state); err != nil {
		return err
	}

	if c.Level == provision.BoxSome {
		log.Debugf("  set state[%s] of container (%s, %s)", c.BoxId, c.Name, state.String())

		if comp, err := carton.NewComponent(c.BoxId, c.AccountId, ""); err != nil {
			return err
		} else if err = comp.SetState(state, c.AccountId, ""); err != nil {
			return err
		}
	}
	return nil
}

func (c *Container) UpdateContId() error {
	log.Debugf("  update Id[%s] of container (%s) in cassandra", c.Id, c.Name)

	var ids = make(map[string][]string)
	cid := []string{c.Id}
	ids["containerId"] = cid
	if asm, err := carton.NewAssembly(c.CartonId, c.AccountId, ""); err != nil {
		return err
	} else if err = asm.NukeAndSetOutputs(ids); err != nil {
		return err
	}
	return nil
}

func urlToHost(urlStr string) string {
	url, _ := url.Parse(urlStr)
	if url == nil || url.Host == "" {
		return urlStr
	}
	host, _, _ := net.SplitHostPort(url.Host)
	if host == "" {
		return url.Host
	}
	return host
}

func (c *Container) StateCheck(args RancherProvisioner) (err error) {
	var res *client.Container
	err = safe.WaitCondition(10*time.Minute, 10*time.Second, func() (bool, error) {
		res, err = args.Cluster().GetContainerById(c.Id)
		if err != nil {
			return false, err
		}
		return (res.State == "running"), nil
	})

	if err != nil {
		return err
	}
	c.HostId = res.HostId
	c.PublicIp = res.PrimaryIpAddress
	return nil
}

type NetworkInfo struct {
	HTTPHostPort string
	IP           string
}

func (c *Container) NetworkInfo(r RancherProvisioner) error {
	err := r.Cluster().SetNetworkinNode(c.HostId, c.PublicIp, c.CartonId, c.AccountId)
	if err != nil {
		return err
	}
	return err
}

/*
func (c *Container) Logs(p DockerProvisioner)   error {
	var outBuffer bytes.Buffer
		var closeChan chan bool
		b := &provision.Box{Id: c.Id, Name: c.BoxName, Tosca: "docker"}
		logWriter := carton.NewLogWriter(b)
		writer := io.MultiWriter(&outBuffer, &logWriter)
	  logopt := docker.LogsOptions{
    Container:  c.Id,
		OutputStream: writer,
		ErrorStream: writer,
			Follow:    true,
			//	RawTerminal:  true,
				Stdout:       true,
				Stderr:       true,
				Timestamps:   false,
			//	Tail:         "100",
		}

	cs := make(chan []byte)
  go p.Cluster().SetLogs(cs,logopt, closeChan)

	go func(closeChan chan bool, logWriter carton.LogWriter) {
		select {
		case <-closeChan:
			logWriter.Close()
		default:
		}
	}(closeChan, logWriter)

var err error
	if err != nil {
		return  err
	}
return  nil

}
func (c *Container) hostToNodeAddress(p DockerProvisioner, host string) (string, error) {
	nodes, err := p.Cluster().Nodes()
	if err != nil {
		return "", err
	}
	for _, node := range nodes {
		if urlToHost(node.Address) == host {
			return node.Address, nil
		}
	}
	return "", fmt.Errorf("Host `%s` not found", host)
}



*/

/*
func (c *Container) addEnvsToConfig(args *CreateArgs, cfg *docker.Config) {
if !args.Deploy {
		for _, envData := range args.Box.Envs() {
			cfg.Env = append(cfg.Env, fmt.Sprintf("%s=%s", envData.Name, envData.Value))
		}
	}

} */
/*
func (c *Container) Remove(p DockerProvisioner) error {
	log.Debugf("Removing container %s from docker", c.BoxName)

	//this will be removed. containerID will be stored upon create in riak
	id, _ := p.Cluster().PreStopAction(c.BoxName)
	c.Id = id
	err := c.Stop(p)
	if err != nil {
		log.Errorf("error on stop unit %s - %s", c.Id, err)
	}
	err = p.Cluster().RemoveContainer(docker.RemoveContainerOptions{ID: c.Id})
	if err != nil {
		log.Errorf("Failed to remove container from docker: %s", err)
	}
	return nil
}

type StartArgs struct {
	Provisioner DockerProvisioner
	Box         *provision.Box
	Deploy      bool
}

func (c *Container) Start(args *StartArgs) error {
	_, err := getPort()
	if err != nil {
		return err
	}

	hostConfig := docker.HostConfig{
		Memory:     int64(args.Box.ConGetMemory()),
		MemorySwap: int64(args.Box.ConGetMemory() + args.Box.GetSwap()),
		CPUShares:  int64(args.Box.GetCpushare()),
	}

	err = args.Provisioner.Cluster().StartContainer(c.Id, &hostConfig)
	if err != nil {
		return err
	}
	initialStatus := constants.StatusContainerStarted
	if args.Deploy {
		initialStatus = constants.StatusContainerBootstrapping
	}
	return c.SetStatus(initialStatus)
}

func (c *Container) Stop(p DockerProvisioner) error {
	if c.Status.String() == constants.StatusContainerStopped.String() {
		return nil
	}
	err := p.Cluster().StopContainer(c.Id, 10)
	if err != nil {
		log.Errorf("error on stop container %s: %s", c.Id, err)
	}
	c.SetStatus(constants.StatusContainerStopped)
	return nil
}

type waitResult struct {
	status int
	err    error
}

var safeAttachInspectTimeout = 20 * time.Second

func SafeAttachWaitContainer(p DockerProvisioner, opts docker.AttachToContainerOptions) (int, error) {
	cluster := p.Cluster()
	resultCh := make(chan waitResult, 1)
	go func() {
		err := cluster.AttachToContainer(opts)
		if err != nil {
			resultCh <- waitResult{err: err}
			return
		}
		status, err := cluster.WaitContainer(opts.Container)
		resultCh <- waitResult{status: status, err: err}
	}()
	for {
		select {
		case result := <-resultCh:
			return result.status, result.err
		case <-time.After(safeAttachInspectTimeout):
		}
		contData, err := cluster.InspectContainer(opts.Container)
		if err != nil {
			return 0, err
		}
		if !contData.State.Running {
			return contData.State.ExitCode, nil
		}
	}
}
*/

/*





type Pty struct {
	Width  int
	Height int
	Term   string
}

func (c *Container) Shell(p DockerProvisioner, stdin io.Reader, stdout, stderr io.Writer, pty Pty) error {
	cmds := []string{"/usr/bin/env", "TERM=" + pty.Term, "bash", "-l"}
	execCreateOpts := docker.CreateExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmds,
		Container:    c.Id,
		Tty:          true,
	}
	exec, err := p.Cluster().CreateExec(execCreateOpts)
	if err != nil {
		return err
	}
	startExecOptions := docker.StartExecOptions{
		InputStream:  stdin,
		OutputStream: stdout,
		ErrorStream:  stderr,
		Tty:          true,
		RawTerminal:  true,
	}
	errs := make(chan error, 1)
	go func() {
		errs <- p.Cluster().StartExec(exec.ID, c.Id, startExecOptions)
	}()
	execInfo, err := p.Cluster().InspectExec(exec.ID, c.Id)
	for !execInfo.Running && err == nil {
		select {
		case startErr := <-errs:
			return startErr
		default:
			execInfo, err = p.Cluster().InspectExec(exec.ID, c.Id)
		}
	}
	if err != nil {
		return err
	}
	p.Cluster().ResizeExecTTY(exec.ID, c.Id, pty.Height, pty.Width)
	return <-errs
}

type execErr struct {
	code int
}

func (e *execErr) Error() string {
	return fmt.Sprintf("unexpected exit code: %d", e.code)
}

func (c *Container) Exec(p DockerProvisioner, stdout, stderr io.Writer, cmd string, args ...string) error {
	cmds := []string{"/bin/bash", "-lc", cmd}
	cmds = append(cmds, args...)
	execCreateOpts := docker.CreateExecOptions{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          cmds,
		Container:    c.Id,
	}
	exec, err := p.Cluster().CreateExec(execCreateOpts)
	if err != nil {
		return err
	}
	startExecOptions := docker.StartExecOptions{
		OutputStream: stdout,
		ErrorStream:  stderr,
	}
	err = p.Cluster().StartExec(exec.ID, c.Id, startExecOptions)
	if err != nil {
		return err
	}
	execData, err := p.Cluster().InspectExec(exec.ID, c.Id)
	if err != nil {
		return err
	}
	if execData.ExitCode != 0 {
		return &execErr{code: execData.ExitCode}
	}
	return nil

}

/*
// Commits commits the container, creating an image in Docker. It then returns
// the image identifier for usage in future container creation.
func (c *Container) Commit(p DockerProvisioner, writer io.Writer) (string, error) {
	log.Debugf("commiting container %s", c.Id)
	parts := strings.Split(c.BuildingImage, ":")
	if len(parts) < 2 {
		return "", log.WrapError(fmt.Errorf("error parsing image name, not enough parts: %s", c.BuildingImage))
	}
	repository := strings.Join(parts[:len(parts)-1], ":")
	tag := parts[len(parts)-1]
	opts := docker.CommitContainerOptions{Container: c.Id, Repository: repository, Tag: tag}
	image, err := p.Cluster().CommitContainer(opts)
	if err != nil {
		return "", log.WrapError(fmt.Errorf("error in commit container %s: %s", c.Id, err.Error()))
	}
	imgData, err := p.Cluster().InspectImage(c.BuildingImage)
	imgSize := ""
	if err == nil {
		imgSize = fmt.Sprintf("(%.02fMB)", float64(imgData.Size)/1024/1024)
	}
	fmt.Fprintf(writer, " ---> Sending image to repository %s\n", imgSize)
	log.Debugf("image %s generated from container %s", image.ID, c.Id)
	maxTry, _ := config.GetInt("docker:registry-max-try")
	if maxTry <= 0 {
		maxTry = 3
	}
	for i := 0; i < maxTry; i++ {
		err = p.PushImage(repository, tag)
		if err != nil {
			fmt.Fprintf(writer, "Could not send image, trying again. Original error: %s\n", err.Error())
			log.Errorf("error in push image %s: %s", c.BuildingImage, err.Error())
			time.Sleep(time.Second)
			continue
		}
		break
	}
	if err != nil {
		return "", log.WrapError(fmt.Errorf("error in push image %s: %s", c.BuildingImage, err.Error()))
	}
	return c.BuildingImage, nil
}
*/
/*
func getPort() (string, error) {
	/*	port, err := config.Get("docker:run-cmd:port")
		if err != nil {
			return "", err
		}
		return fmt.Sprint(port), nil
*/
//	return "", nil
//}
