// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package repository

func init() {
	Register("nop", nopManager{})
}

type nopManager struct{}

func (nopManager) CreateHook(username string) error {
	return nil
}

func (nopManager) RemoveHook(username string) error {
	return nil
}
