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
	"strings"

)

type boxImages struct {
	BoxName string
	Images  []string
	Count   int
}

type ImageMetadata struct {
	Name string
}

func boxImageName(appName string) string {
	return fmt.Sprintf("%s/app-%s", appName)
}

func listBoxImages(appName string) ([]string, error) {
	//call marketplaces bucket
	//we need to improve db.FetchCollection code to list all from the bucket.
	return []string{}, nil
}

//list all the tosca images from marketplace and figure out if this is valid.
func isValidBoxImage(appName, imageId string) (bool, error) {
	images, err := listBoxImages(appName)

	if err != nil {
		return false, err
	}
	for _, img := range images {
		if img == imageId {
			return true, nil
		}
	}
	//return false, nil
	return true, nil

}
