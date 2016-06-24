/*
** Copyright [2013-2016] [Megam Systems]
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
package provision

import (
	"regexp"
	"strconv"
	"strings"
	"fmt"

	"github.com/pivotal-golang/bytefmt"
)

type BoxCompute struct {
	Cpushare string
	Memory   string
	Swap     string
	HDD      string
}

func (bc *BoxCompute) trimCore() string {
	coreRegex := regexp.MustCompile("[Cc][Oo][Rr][Ee].*")
	return strings.TrimSpace(coreRegex.ReplaceAllString(bc.Cpushare, ""))
}
func (bc *BoxCompute) trimMemory() string {
	coreRegex := regexp.MustCompile("[Bb].*")
	s := strings.TrimSpace(coreRegex.ReplaceAllString(bc.Memory, ""))
	return  strings.Replace(s, " ", "", -1)
}


func (bc *BoxCompute) numCpushare() uint64 {
	if cs, err := strconv.ParseUint(bc.trimCore(), 10, 64); err != nil {
		return 0
	} else {
		return cs
	}
}

func (bc *BoxCompute) ConnumMemory() uint64 {
	if cs, err := strconv.ParseUint(bc.trimMemory(), 10, 64); err != nil {
		return 0
	} else {
		return cs
	}
}

func (bc *BoxCompute) numMemory() uint64 {
	if cp, err := bytefmt.ToMegabytes(strings.Replace(bc.Memory, " ", "", -1)); err != nil {
		return 0
	} else {
		return cp
	}
}

func (bc *BoxCompute) numSwap() uint64 {
	if cs, err := bytefmt.ToMegabytes(bc.Swap); err != nil {
		return 0
	} else {
		return cs
	}
}

func (bc *BoxCompute) trimHDD() string {
	sataRegex := regexp.MustCompile("[Ss][aA][Tt][aA].*")
	ssdRegex := regexp.MustCompile("[Ss][sS][dD].*")
	var hddTrim = ""
	if len(sataRegex.FindStringSubmatch(bc.HDD)) > 0 {
		hddTrim = strings.TrimSpace(sataRegex.ReplaceAllString(bc.HDD, ""))
	}

	if len(ssdRegex.FindStringSubmatch(bc.HDD)) > 0 {
		hddTrim = strings.TrimSpace(ssdRegex.ReplaceAllString(bc.HDD, ""))
	}
	return strings.Replace(hddTrim, " ", "", -1)
}

func (bc *BoxCompute) numHDD() uint64 {
	if cp, err := bytefmt.ToMegabytes(bc.trimHDD()); err != nil {
		return 10240
	} else {
		return cp
	}
}

func (bc *BoxCompute) String() string {
	return "(" + strings.Join([]string{
		CPU + ":" + strconv.FormatInt(int64(bc.numCpushare()), 10),
		RAM + ":" + strconv.FormatInt(int64(bc.numMemory()), 10),
		HDD + ":" + strconv.FormatInt(int64(bc.numHDD()), 10)},
		",") + " )"
}
