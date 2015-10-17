package db

import (
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/db"
	"github.com/megamsys/libgo/hc"

	"github.com/megamsys/megamd/meta"
)

func init() {
	hc.AddChecker("Riak", healthCheck)
}

func healthCheck() error {
	conn, err := newConn("test")
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

//A global function which helps to avoid passing config of riak everywhere.
func newConn(bkt string) (*db.Storage, error) {
	r, err := db.NewRiakDB(meta.MC.Riak, bkt)

	if err != nil {
		return nil, err
	}

	return r.Conn()
}

func Fetch(bkt string, key string, data interface{}) error {
	s, err := newConn(bkt)
	if err != nil {
		return err
	}
	defer s.Close()
	log.Debugf("%s (%s, %s)", cmd.Colorfy("  > [riak] fetch", "blue", "", "bold"), bkt, key)

	if err = s.FetchStruct(key, data); err != nil {
		return err
	}
	return nil
}

func Store(bkt string, key string, data interface{}) error {
	s, err := newConn(bkt)
	if err != nil {
		return err
	}
	defer s.Close()
	log.Debugf("%s (%s, %s)", cmd.Colorfy("  > [riak] store", "blue", "", "bold"), bkt, key)

	if err = s.StoreStruct(key, data); err != nil {
		return err
	}
	return nil
}
