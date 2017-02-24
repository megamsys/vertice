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
package marketplaces

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
	"time"
)

var (
	// categories of marketplaces queue process
	RAWIMAGE     = "rawimage"
	MARKETPLACES = "marketplaces"
)

type Requests struct {
	Id        string    `json:"id" cql:"id"`
	Name      string    `json:"name" cql:"name"`
	AccountId string    `json:"account_id" cql:"account_id"`
	CatId     string    `json:"cat_id" cql:"cat_id"`
	Action    string    `json:"action" cql:"action"`
	Category  string    `json:"category" cql:"category"`
	CreatedAt time.Time `json:"created_at" cql:"created_at"`
}

type ApiRequests struct {
	JsonClaz string     `json:"json_claz" cql:"json_claz"`
	Results  []Requests `json:"results" cql:"results"`
}

func (r *Requests) String() string {
	if d, err := yaml.Marshal(r); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

func (r *Requests) ParseRequest() (MarketplaceInterface, error) {
	switch r.Category {
	case RAWIMAGE:
		return r.getRawImage()
	case MARKETPLACES:
		return r.getMarketplace()
	default:
		return nil, newParseError([]string{r.Category, r.Action}, []string{RAWIMAGE, MARKETPLACES})
	}
}

func (r *Requests) getRawImage() (*RawImages, error) {
	raw := new(RawImages)
	raw.AccountId = r.AccountId
	raw.Id = r.CatId
	return raw.Get()
}

func (r *Requests) getMarketplace() (*Marketplaces, error) {
	return GetMarketplace(r.AccountId, r.CatId)
}

// ParseError represents an error that occurred during parsing.
type ParseError struct {
	Found    string
	Expected []string
}

// newParseError returns a new instance of ParseError.
func newParseError(found []string, expected []string) *ParseError {
	return &ParseError{Found: strings.Join(found, ","), Expected: expected}
}

// Error returns the string representation of the error.
func (e *ParseError) Error() string {
	return fmt.Sprintf("found %s, expected %s", e.Found, strings.Join(e.Expected, ", "))
}
