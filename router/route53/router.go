// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package route53

import (
	"crypto/md5"
	"fmt"
	"net/url"

	"github.com/megamsys/libgo"
	"github.com/megamsys/megamd/router"
)

const routerName = "route53"

type route53Router struct {
	client *api.Client
	prefix string
	domain string
}

func createRouter(prefix string) (router.Router, error) {
	vURL, err := config.GetString(prefix + ":api-url")
	if err != nil {
		return nil, err
	}
	domain, err := config.GetString(prefix + ":domain")
	if err != nil {
		return nil, err
	}
	client := api.NewClient(vURL, registry.GetRegistry())
	vRouter := &route53Router{
		client: client,
		prefix: prefix,
		domain: domain,
	}
	return vRouter, nil
}


func (r *route53Router) AddRoute(name string, address *url.URL) error {
	usedName, err := router.Retrieve(name)
	if err != nil {
		return err
	}
	serverKey := engine.ServerKey{
		Id:         r.serverName(address.String()),
		BackendKey: engine.BackendKey{Id: r.backendName(usedName)},
	}
	if found, _ := r.client.GetServer(serverKey); found != nil {
		return router.ErrRouteExists
	}
	server, err := engine.NewServer(serverKey.Id, address.String())
	if err != nil {
		return &router.RouterError{Err: err, Op: "add-route"}
	}
	err = r.client.UpsertServer(serverKey.BackendKey, *server, engine.NoTTL)
	if err != nil {
		return &router.RouterError{Err: err, Op: "add-route"}
	}
	return nil
}

func (r *route53Router) RemoveRoute(name string, address *url.URL) error {
	usedName, err := router.Retrieve(name)
	if err != nil {
		return err
	}
	serverKey := engine.ServerKey{
		Id:         r.serverName(address.String()),
		BackendKey: engine.BackendKey{Id: r.backendName(usedName)},
	}
	err = r.client.DeleteServer(serverKey)
	if err != nil {
		if _, ok := err.(*engine.NotFoundError); ok {
			return router.ErrRouteNotFound
		}
		return &router.RouterError{Err: err, Op: "remove-route"}
	}
	return nil
}

func (r *route53Router) SetCName(cname, name string) error {
	usedName, err := router.Retrieve(name)
	if err != nil {
		return err
	}
	if !router.ValidCName(cname, r.domain) {
		return router.ErrCNameNotAllowed
	}
	frontendName := r.frontendName(cname)
	if found, _ := r.client.GetFrontend(engine.FrontendKey{Id: frontendName}); found != nil {
		return router.ErrCNameExists
	}
	frontend, err := engine.NewHTTPFrontend(
		frontendName,
		r.backendName(usedName),
		fmt.Sprintf(`Host(%q)`, cname),
		engine.HTTPFrontendSettings{},
	)
	if err != nil {
		return &router.RouterError{Err: err, Op: "set-cname"}
	}
	err = r.client.UpsertFrontend(*frontend, engine.NoTTL)
	if err != nil {
		return &router.RouterError{Err: err, Op: "set-cname"}
	}
	return nil
}

func (r *route53Router) UnsetCName(cname, _ string) error {
	frontendKey := engine.FrontendKey{Id: r.frontendName(cname)}
	err := r.client.DeleteFrontend(frontendKey)
	if err != nil {
		if _, ok := err.(*engine.NotFoundError); ok {
			return router.ErrCNameNotFound
		}
		return &router.RouterError{Err: err, Op: "unset-cname"}
	}
	return nil
}

func (r *route53Router) Addr(name string) (string, error) {
	usedName, err := router.Retrieve(name)
	if err != nil {
		return "", err
	}
	frontendHostname := r.frontendHostname(usedName)
	frontendKey := engine.FrontendKey{Id: r.frontendName(frontendHostname)}
	if found, _ := r.client.GetFrontend(frontendKey); found == nil {
		return "", router.ErrRouteNotFound
	}
	return frontendHostname, nil
}


func (r *route53Router) StartupMessage() (string, error) {
	message := fmt.Sprintf("route53 router %q with API at %q", r.domain, r.client.Addr)
	return message, nil
}
