// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package router provides interfaces that need to be satisfied in order to
// implement a new router on tsuru.
package router

import (
	"errors"
	"fmt"
)

type routerFactory func(string) (Router, error)

var (
	ErrRouteExists     = errors.New("Route already exists")
	ErrRouteNotFound   = errors.New("Route not found")
	ErrCNameExists     = errors.New("CName already exists")
	ErrCNameNotFound   = errors.New("CName not found")
	ErrCNameNotAllowed = errors.New("CName as router subdomain not allowed")
)

var routers = make(map[string]routerFactory)

// Register registers a new router.
func Register(name string, r routerFactory) {
	routers[name] = r
}

// Get gets the named router from the registry.
func Get(name string) (Router, error) {
	/*prefix := "routers:" + name
	routerType, err := config.GetString(prefix + ":type")
	if err != nil {
		msg := fmt.Sprintf("config key '%s:type' not found", prefix)
		if name != "hipache" {
			return nil, errors.New(msg)
		}
		log.Errorf("WARNING: %s, fallback to top level '%s:*' router config", msg, name)
		routerType = name
		prefix = name
	}
	factory, ok := routers[routerType]
	if !ok {
		return nil, fmt.Errorf("unknown router: %q.", routerType)
	}
	r, err := factory(prefix)
	if err != nil {
		return nil, err
	}
	return r, nil
	*/
	return nil, nil
}

// Router is the basic interface of this package. It provides methods for
// managing backends and routes. Each backend can have multiple routes.
type Router interface {
	SetCName(cname, name string) error
	UnsetCName(cname, name string) error
	Addr(name string) (string, error)
}

type MessageRouter interface {
	StartupMessage() (string, error)
}

type RouterError struct {
	Op  string
	Err error
}

func (e *RouterError) Error() string {
	return fmt.Sprintf("[router %s] %s", e.Op, e.Err)
}


/* Store stores the app name related with the
// router name.
func Store(appName, routerName, kind string) error {
	coll, err := collection()
	if err != nil {
		return err
	}
	defer coll.Close()
	data := map[string]string{
		"app":    appName,
		"router": routerName,
		"kind":   kind,
	}
	return coll.Insert(&data)
}

func retrieveRouterData(appName string) (map[string]string, error) {
	data := map[string]string{}
	coll, err := collection()
	if err != nil {
		return data, err
	}
	defer coll.Close()
	err = coll.Find(bson.M{"app": appName}).One(&data)
	// Avoid need for data migrations, before kind existed we only supported
	// hipache as a router so we set is as default here.
	if data["kind"] == "" {
		data["kind"] = "hipache"
	}
	return data, err
}

func Retrieve(appName string) (string, error) {
	data, err := retrieveRouterData(appName)
	if err != nil {
		if err == mgo.ErrNotFound {
			return "", ErrBackendNotFound
		}
		return "", err
	}
	return data["router"], nil
}

func Remove(appName string) error {
	coll, err := collection()
	if err != nil {
		return err
	}
	defer coll.Close()
	return coll.Remove(bson.M{"app": appName})
}

// validCName returns true if the cname is not a subdomain of
// the router current domain, false otherwise.
func ValidCName(cname, domain string) bool {
	return !strings.HasSuffix(cname, domain)
}
*/
