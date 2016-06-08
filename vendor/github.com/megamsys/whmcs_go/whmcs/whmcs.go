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

package whmcs

import (
	//"encoding/json"
	"errors"
	"fmt"
	//"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	apiEndSux = "includes/api.php"

	ACTION       = "action"
	ACCESSKEY    = "accesskey"
	RESPONSETYPE = "json"
	USERNAME     = "username"
	PASSWORD     = "password"
)

var (
	ErrAuthenticationKeysNotFound = errors.New("[username, password] or access_key not found. Did you send it ?")
	ErrActionNotFound             = errors.New("[action] not found. Did you send it ?")
)

// A Client manages communication with the WHMCS API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.  Defaults to https://www.yourdomain.com/billing/, but can be
	// set to a domain endpoint to use with your billing at your enterprise.  BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL

	// API endpoint suffix used when communicating with the WHMCS API.
	ApiEndSux string

	// Services used for talking to different parts of the WHMCS API.
	Orders    *OrdersService
	Billables *BillablesService
	Accounts  *AccountsService
}

// A WRequest manages communication with the WHMCS API.
type WRequest struct {
	data *url.Values

	url *url.URL
}

func (w WRequest) howzzat(key string) bool {
	return len(w.data.Get(key)) > 0
}

type Params struct {
	parms map[string]string
	u     string
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}

// addFormValues adds the parameters in opt as URL values parameters.
func addFormValues(opt map[string]string) *url.Values {
	uv := url.Values{}
	for k, v := range opt {
		uv.Set(k, v)
	}
	return &uv
}

// NewClient returns a new WHMCS API client.  If a nil httpClient is
// provided, http.DefaultClient will be used.  To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
func NewClient(httpClient *http.Client, defaultBaseURL string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, ApiEndSux: apiEndSux}

	c.Orders = &OrdersService{client: c}
	c.Billables = &BillablesService{client: c}
	c.Accounts = &AccountsService{client: c}
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewWRequest(dat map[string]string, action string) (*WRequest, error) {
	rel, err := url.Parse(apiEndSux)
	if err != nil {
		return nil, err
	}
	u := c.BaseURL.ResolveReference(rel)

	if len(strings.TrimSpace(action)) > 0 {
		dat[ACTION] = action
		dat[RESPONSETYPE] = "json"
	}
	return &WRequest{url: u, data: addFormValues(dat)}, nil
}

// Response is a WHMCS API response.  This wraps the standard http.Response
// returned from WHMCS and provides convenient access to things like
// pagination links.
type Response struct {
	Status     string // e.g. "200 OK"
    StatusCode int    // e.g. 200
	Body string
	ContentLength int64
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	body, _ := ioutil.ReadAll(r.Body)
	response := &Response{
		Status: r.Status,
		StatusCode: r.StatusCode,
		Body: string(body),
		ContentLength: r.ContentLength,
		}
	return response
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req WRequest, v interface{}) (*Response, error) {
	url, err := c.auth(req)
	if err != nil {
		return nil, err
	}
 fmt.Println("--- " + url.String())
	resp, err := c.client.PostForm(url.String()+"?accesskey=team4megam", *req.data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)
	
	err = CheckResponse(response)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	/*if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, response.Body)
		} else {
			err = json.NewDecoder(response.Body).Decode(v)          
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}*/
	return response, err
}

func do(c *Client, p Params, a interface{}) (*Response, error) {
	req, err := c.NewWRequest(p.parms, p.u)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(*req, a)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func (c *Client) auth(w WRequest) (*url.URL, error) {
	if !((w.howzzat(USERNAME) && w.howzzat(PASSWORD)) || (w.howzzat(ACCESSKEY))) {
		return nil, ErrAuthenticationKeysNotFound
	}

	if !(w.howzzat(ACTION)) {
		return nil, ErrActionNotFound
	}
	//patch if you sent access_key
	if w.howzzat(ACCESSKEY) {
		w.url.Query().Set(ACCESSKEY, w.data.Get(ACCESSKEY))
	}
	return w.url, nil
}

/*
An ErrorResponse reports one or more errors caused by an API request.

WHMCS API docs: http://developer.WHMCS.com/v3/#client-errors
*/
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
	Errors   []Error        `json:"errors"`  // more detail on individual errors
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, sanitizeURL(r.Response.Request.URL),
		r.Response.StatusCode, r.Message, r.Errors)
}

// sanitizeURL redacts the client_secret parameter from the URL which may be
// exposed to the user, specifically in the ErrorResponse error message.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get(ACCESSKEY)) > 0 {
		params.Set(ACCESSKEY, "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}

/*
An Error reports more details on an individual error in an ErrorResponse.
These are the possible validation error codes:

    missing:
        resource does not exist
    missing_field:
        a required field on a resource has not been set
    invalid:
        the formatting of a field is invalid
    already_exists:
        another resource has the same valid as this field

WHMCS API docs: http://developer.WHMCS.com/v3/#client-errors
*/
type Error struct {
	Resource string `json:"resource"` // resource on which the error occurred
	Field    string `json:"field"`    // field on which the error occurred
	Code     string `json:"code"`     // validation error code
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v error caused by %v field on %v resource",
		e.Code, e.Field, e.Resource)
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(r *Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	//errorResponse := &ErrorResponse{Response: r}
	//data, err := ioutil.ReadAll(r.Body)
	//if err == nil && data != nil {
	//	json.Unmarshal(r.Body, errorResponse)
	//}
	return errors.New(r.Body)
}

// parseBoolResponse determines the boolean result from a WHMCS API response.
// Several WHMCS API methods return boolean responses indicated by the HTTP
// status code in the response (true indicated by a 204, false indicated by a
// 404).  This helper function will determine that result and hide the 404
// error if present.  Any other error will be returned through as-is.
func parseBoolResponse(err error) (bool, error) {
	if err == nil {
		return true, nil
	}

	if err, ok := err.(*ErrorResponse); ok && err.Response.StatusCode == http.StatusNotFound {
		// Simply false.  In this one case, we do not pass the error through.
		return false, nil
	}

	// some other real error occurred
	return false, err
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// Int is a helper routine that allocates a new int32 value
// to store v and returns a pointer to it, but unlike Int32
// its argument value is an int.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}
