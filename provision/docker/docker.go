package docker

import (
	"crypto"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"net/url"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/megamsys/megamd/provision/docker/cluster"
)

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

func (p *dockerProvisioner) hostToNodeAddress(host string) (string, error) {
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

// PushImage sends the given image to the registry server defined in the
// configuration file.
func (p *dockerProvisioner) PushImage(name, tag string) error {
	registry := "pull it from box.Repo.Registry.ServerAddress"
	var buf safe.Buffer
	pushOpts := docker.PushImageOptions{Name: name, Tag: tag, OutputStream: &buf}
	err = p.Cluster().PushImage(pushOpts, p.RegistryAuthConfig())
	if err != nil {
		log.Errorf("[docker] Failed to push image %q (%s): %s", name, err, buf.String())
		return err
	}
	return nil
}

func (p *dockerProvisioner) RegistryAuthConfig() docker.AuthConfiguration {
	var authConfig docker.AuthConfiguration
	authConfig.Email = "pull it from box.Repo.Email"
	authConfig.Username = "pull it from box.Repo.Registry.USername"
	authConfig.Password = "pull it from box.Repo.Registry.Password"
	authConfig.ServerAddress = "pull it from box.Repo.Registry.Serveraddress"
	return authConfig
}
