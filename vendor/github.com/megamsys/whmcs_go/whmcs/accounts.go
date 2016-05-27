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

// AccountsService handles communication with the user related
// methods of the WHMCS API.
//
// WHMCS API docs: http://docs.whmcs.com/API#Client_Management
type AccountsService struct {
	client *Client
}

// Account represents an WHMCS user.
type Account struct {
	FirstName   *string `json:"firstname"`
	LastName    *string `json:"lastname"`
	Email       *string `json:"email"`
	Address1    *string `json:"address1"`
	City        *string `json:"city"`
	State       *string `json:"state"`
	PostCode    *string `json:"postcode"`
	Country     *string `json:"country"`
	PhoneNumber *string `json:"phonenumber"`
	Password    *string `json:"password2"`
	Status      *string `json:"status"`
}

func (u Account) String() string {
	return Stringify(u)
}

// Create adds a new user onboarded in Megam vertice
//
// WHMCs API docs: http://docs.whmcs.com/API:Add_Client
func (s *AccountsService) Create(parms map[string]string) (*Account, *Response, error) {
	a := new(Account)
	resp, err := do(s.client, Params{parms: parms, u: "addclient"}, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Get fetches a user.  Passing the empty string will fetch the authenticated
// user.
//
// WHMCS API docs: http://docs.whmcs.com/API:Get_Clients_Details
func (s *AccountsService) Get(parms map[string]string) (*Account, *Response, error) {
	a := new(Account)
	resp, err := do(s.client, Params{parms: parms, u: "getclientsdetails"}, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}

// Edit the authenticated user.
//
// WHMCS API docs: http://docs.whmcs.com/API:Update_Client
func (s *AccountsService) Edit(parms map[string]string) (*Account, *Response, error) {
	a := new(Account)
	resp, err := do(s.client, Params{parms: parms, u: "updateclient"}, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}
