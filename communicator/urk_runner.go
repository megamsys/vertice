package communicator

import (
	"io"
	"os"
	"fmt"
	"bytes"
	"reflect"
	"github.com/megamsys/megdc/handler"
)


var mainurkRunner *UrkRunner

func init() {
	mainurkRunner = &UrkRunner{}
	Register("urknall", mainurkRunner)
}

type  UrkRunner struct{
 wrappedparms *handler.WrappedParms
 OutBuffer bytes.Buffer
 Inputs  map[string]string
}

func NewUrkRunner(i interface{},inputs map[string]string) (Runner, error) {
	ur := new(UrkRunner)
	ur.wrappedparms = handler.NewWrap(i)
	ur.Inputs = inputs //map[string]string{constants.USERMAIL : b.UserMail, constants.HOST_ID: b.HostId} // Email is needed for urknall template's events trigger
  return ur, nil
}

func (ur UrkRunner) Run(packages []string) (r Runner, err error) {
	var writer io.Writer
	writer = io.MultiWriter(&ur.OutBuffer, os.Stdout)
	ur.wrappedparms.IfNoneAddPackages(packages)
	if h, err := handler.NewHandler(ur.wrappedparms); err != nil {
		return nil, err
	} else if err = h.Run(writer, ur.Inputs); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return
}

func (runner UrkRunner) Rerun() (r Runner, err error) {
	return
}

func (runner UrkRunner) CleanUp() (err error) {
	return
}

func (runner UrkRunner) String() string {
  return ""
}

func Clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}
