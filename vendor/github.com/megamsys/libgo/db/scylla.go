package db

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
)

type Options struct {
	TableName string
	Pks       []string
	Ccms      []string
	Username  string
	Password  string
	Keyspace  string
	Hosts     []string
	PksClauses   map[string]interface{}
	CcmsClauses  map[string]interface{}
}

//A global function which helps to avoid passing config of riak everywhere.
func newDBConn(ops Options) (*ScyllaDB, error) {
	r, err := newScyllaDB(ScyllaDBOpts{
		KeySpaceName: ops.Keyspace,
		NodeIps:      ops.Hosts,
		Username:     ops.Username,
		Password:     ops.Password,
		Debug:        true,
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func newScyllaTable(ops Options, data interface{}) *ScyllaTable {
	t, err := newDBConn(ops)
	if err != nil {
		return nil
	}
	log.Debugf("%s (%s, %s)", cmd.Colorfy("  > [scylla] fetch", "blue", "", "bold"), ops.TableName)
	tbl := t.table(ops.TableName, ops.Pks, ops.Ccms, data)
	return tbl
}

func Fetchdb(tinfo Options, data interface{}) error {
	d := newScyllaTable(tinfo, data)
	if d != nil {
		//err := d.read(ScyllaWhere{Clauses: tinfo.Clauses}, data)
		err := d.read(tinfo.PksClauses, tinfo.CcmsClauses, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func Storedb(tinfo Options, data interface{}) error {
	t := newScyllaTable(tinfo, data)
	err := t.insert(data)
	if err != nil {
		return err
	}
	return nil
}

func Updatedb(tinfo Options, data map[string]interface{}) error {
	t := newScyllaTable(tinfo, data)
	err := t.update(tinfo, data)
	if err != nil {
		return err
	}
	return nil
}

func Deletedb(tinfo Options, data interface{}) error {
	t := newScyllaTable(tinfo, data)
	err := t.deleterow(tinfo)
	if err != nil {
		return err
	}
	return nil
}


