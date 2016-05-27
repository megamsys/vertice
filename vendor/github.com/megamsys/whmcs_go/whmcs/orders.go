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

// OrdersService provides access to the orders related functions
// in the WHMCS API.
//
// WHMCS API docs: http://docs.whmcs.com/API
type OrdersService struct {
	client *Client
}

// Orders represents a WHCMS Order for an account
type Order struct {
	ClientId      *string `json:"clientid"`
	PID           *string `json:"pid"`
	Domain        *string `json:"domain`
	BillingCycle  *string `json:"billingcycle"`
	DomainType    *string `json:"domaintype"`
	RegPeriod     *string `json:"regperiod"`
	EppCode       *int    `json:"eppcode"`
	NameServer1   *string `json:"nameserver1"`
	PaymentMethod *string `json:"paymentmethod"`
	HostName      *string `json:"hostname"`
}

func (o Order) String() string {
	return Stringify(o)
}

// Create adds an  new order
//
// WHMCs API docs: http://docs.whmcs.com/API:Add_Order
func (s *OrdersService) Create(parms map[string]string) (*Order, *Response, error) {
	order := new(Order)
	resp, err := do(s.client, Params{parms: parms, u: "addorder"}, order)
	if err != nil {
		return nil, resp, err
	}
	return order, resp, err
}

// List the orders for a user.  Passing the empty string will list
// orders for the authenticated user.
//
// WHMCS API docs: http://docs.whmcs.com/API:Get_Orders
func (s *OrdersService) List(parms map[string]string) (*[]Order, *Response, error) {
	orders := new([]Order)
	resp, err := do(s.client, Params{parms: parms, u: "getorders"}, orders)
	if err != nil {
		return nil, resp, err
	}
	return orders, resp, err
}

// Get the status of the order
//
// WHMCS API docs: http://docs.whmcs.com/API:Get_Order_Statuses
// TO-DO this shall return *[]OrderStatus
func (s *OrdersService) Status(parms map[string]string) (*Order, *Response, error) {
	order := new(Order)
	resp, err := do(s.client, Params{parms: parms, u: "getorderstatuses"}, order)
	if err != nil {
		return nil, resp, err
	}
	return order, resp, err
}

// Cancel an order.
//
// WHMCS API docs: http://docs.whmcs.com/API:Cancel_Order
func (s *OrdersService) Cancel(parms map[string]string) (*Order, *Response, error) {
	order := new(Order)
	resp, err := do(s.client, Params{parms: parms, u: "cancelorder"}, order)
	if err != nil {
		return nil, resp, err
	}
	return order, resp, err
}
