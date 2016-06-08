package bills

import (
	"reflect"
)

const (
	SCYLLADB = "scylladb"
	WHMCS    = "whmcs"
)

var BillProviders map[string]BillProvider

//BillOpts represents a billtransaction managed by the provider
type BillOpts struct {
	AccountId  string
	AssemblyId string
	AssemblyName string
	Consumed   string
	StartTime  string
	EndTime    string
	Timestamp  string
	M   map[string]string
}

type BillProvider interface {
	IsEnabled() bool               //is this billing provider enabled.
	Onboard(o *BillOpts, m map[string]string) error     //onboard an user in the billing system
	Nuke(o *BillOpts) error        //delete an user from the billing system
	Suspend(o *BillOpts) error     //suspend an user from the billing system
	Deduct(o *BillOpts, m map[string]string) error      //deduct the balance credit
	Transaction(o *BillOpts, m map[string]string) error //deduct the bill transaction
	Invoice(o *BillOpts) error     //invoice for a  range.
	Notify(o *BillOpts) error      //notify an user
}

// Provider returns the current configured manager, as defined in the
// configuration file.
func Provider(providerName string) BillProvider {
	if _, ok := BillProviders[providerName]; !ok {
		providerName = "nop"
	}
	return BillProviders[providerName]
}

// Register registers a new billing provider, that can be later configured
// and used.
func Register(name string, provider BillProvider) {
	if BillProviders == nil {
		BillProviders = make(map[string]BillProvider)
	}
	BillProviders[name] = provider
}

func SetField(obj interface{}, name string, value string) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	//structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	//if structFieldType != val.Type() {
		//invalidTypeError := errors.New("Provided value type didn't match obj field type")
	//}

	structFieldValue.Set(val)
	return nil
}

func (s *BillOpts) FillStruct(m map[string]string) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
