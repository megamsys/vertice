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
package carton

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
)

var (
	STATE   = "state"
	CONTROL = "control"
	POLICY  = "policy"

	CREATE    = "create"
	DELETE    = "delete"
	STOP      = "stop"
	START     = "start"
	RESTART   = "restart"
	STATEUP   = "stateup"
	STATEDOWN = "statedown"
	BIND      = "bind"
	UNBIND    = "unbind"
)

type ReqParser struct {
}

// NewParser returns a new instance of Parser.
func NewReqParser(r io.Reader) *ReqParser {
	return &ReqParser{}
}

// ParseStatement parses a statement string and returns its AST representation.
func ParseRequest(s string) (MegdProcessor, error) {
	return NewReqParser(strings.NewReader(s)).ParseRequest(s, "")
}

func (p *ReqParser) ParseRequest(category string, action string) (MegdProcessor, error) {
	switch category {
	case STATE:
		return p.parseState(action)
	case CONTROL:
		return p.parseControl(action)
	case POLICY:
		return p.parsePolicy(action)
	default:
		return nil, newParseError([]string{category, action}, []string{STATE, CONTROL, POLICY})
	}
}

func (p *ReqParser) parseState(action string) (MegdProcessor, error) {
	switch action {
	case CREATE:
		return CreateProcess{}, nil
	case STATEUP:
		return StateupProcess{}, nil
	case STATEDOWN:
		return StateupProcess{}, nil
	default:
		return nil, newParseError([]string{STATE, action}, []string{STATEUP, STATEDOWN})
	}
}

func (p *ReqParser) parseControl(action string) (MegdProcessor, error) {
	switch action {
	case START:
		return StartProcess{}, nil
	case STOP:
		return StopProcess{}, nil
	case RESTART:
		return RestartProcess{}, nil
	default:
		return nil, newParseError([]string{CONTROL, action}, []string{STATEUP, STATEDOWN})
	}
}

func (p *ReqParser) parsePolicy(action string) (MegdProcessor, error) {
	switch action {
	case BIND:
		//	return BindPolicy{}
		return StartProcess{}, nil
	case UNBIND:
		return StopProcess{}, nil
	//	return UnBindPolicy{}
	default:
		return nil, newParseError([]string{POLICY, action}, []string{STATEUP, STATEDOWN})
	}
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

type Requests struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CatId     string `json:"cat_id"`
	CatType   string `json:"cattype"`
	CreatedAt string `json:"created_at"`
}

func (r *Requests) String() string {
	if d, err := yaml.Marshal(r); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}
