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

// BillablesService handles communication with the billableitem related
// methods of the WHMCS API.
//
// WHMCS API docs: http://docs.whmcs.com/API
type BillablesService struct {
	client *Client
}

// BillableItem represents a GitHub repository.
type BillableItem struct {
	ClientId      *string `json:"clientid"`
	Description   *string `json:"description"`
	Hours         *string `json:"hours"`
	Amount        *string `json:"amount"`
	InvoiceAction *string `json:"invoiceaction"`
}

func (r BillableItem) String() string {
	return Stringify(r)
}

// Create a new billable item.
//
// WHMCS API docs: http://docs.whmcs.com/API:Add_Billable_Item
func (s *BillablesService) Create(parms map[string]string) (*BillableItem, *Response, error) {
	a := new(BillableItem)
	resp, err := do(s.client, Params{parms: parms, u: "addbillableitem"}, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, err
}
