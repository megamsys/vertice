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
package cmd

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

const defaultTOSCAPath = "conf/tosca_schema.yaml"

func NewTOSCA() string {
	p, _ := filepath.Abs(defaultTOSCAPath)
	log.Println(fmt.Errorf("Conf: %s", p))

	data, err := ioutil.ReadFile(p)

	if err != nil {
		log.Fatalf("error: %v", err)
	}

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for key, value := range m {
		log.Printf("\n%v\n :=>\n", key)
		switch value.(type) {
		case string:
			log.Printf("===[%v]\n", value)
		case map[interface{}]interface{}:
			log.Printf("???[%v]\n", value)
			break
		default:
			log.Printf(">>>[%v] is unknown!!\n", value)
		}
	}

	d, err := yaml.Marshal(&m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//  log.Printf("--- m dump:\n%s\n\n", string(d))
	return string(d)
}
