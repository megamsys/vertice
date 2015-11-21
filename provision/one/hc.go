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
	//"github.com/megamsys/provision/one/cluster"
	//"github.com/megamsys/opennebula-go/api"
	"github.com/megamsys/libgo/hc"
)


func init() {
	hc.AddChecker("one", healthCheck)
}

func healthCheck() error {
	/*
  we need to pass the Onedeployd config.
	var nodes []cluster.Node = []cluster.Node{cluster.Node{
		Address:  m[api.ENDPOINT],
		Metadata: m,
	},
	}
	cluster, err = cluster.New(&cluster.MapStorage{}, nodes...)
	nodlist, err := c.Nodes()

  if err != nil || len(nodlist) <= 0 {
 	 return err
  }*/
	return nil
}
