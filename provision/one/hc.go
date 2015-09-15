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

package one

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/megamsys/libgo/hc"
)

var httpRegexp = regexp.MustCompile(`^https?://`)

func init() {
	hc.AddChecker("one", healthCheckOne)
}

func healthCheckOne() error {
	onerpc := "http://192.168.1.100/xmlrpc" //We need to hookup deployd.OneEndPoint
	if onerpc == "" {
		return hc.ErrDisabledComponent
	}
	if !httpRegexp.MatchString(onerpc) {
		onerpc = "http://" + onerpc
	}
	onerpc = strings.TrimRight(onerpc, "/")
	resp, err := http.Get(onerpc)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status - %s", body)
	}
	return nil
}
