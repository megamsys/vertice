package alerts

import (
     ldb "github.com/megamsys/libgo/db"
     log "github.com/Sirupsen/logrus"
     constants "github.com/megamsys/libgo/utils"
	"time"
	"strings"
	)

const EVENTSBUCKET = "events"

type Scylla struct {
	scylla_host         []string 
	scylla_keyspace     string   
}

type Events struct {
	EventType   string   	`json:"event_type" cql:"event_type"`
	AccountId   string		`json:"account_id" cql:"account_id"`
	AssemblyId  string		`json:"assembly_id" cql:"assembly_id"`
	Data         []string	`json:"data" cql:"data"`
	CreatedAt   string		`json:"created_at" cql:"created_at"`
}

func NewScylla(m map[string]string) Notifier {
	return &Scylla{
		scylla_host:  strings.Split(m[constants.SCYLLAHOST], ","),
		scylla_keyspace: m[constants.SCYLLAKEYSPACE],
		}
}

func (s *Scylla) satisfied(eva EventAction) bool {
	if eva == STATUS {
		return true
	}
	return false
}


func (s *Scylla) Notify(eva EventAction, edata EventData) error {
	if !s.satisfied(eva) {
		return nil
	}
	s_data := parseMapToOutputFormat(edata)
	ops := ldb.Options{
			TableName:   EVENTSBUCKET,
			Pks:         []string{constants.EVENT_TYPE, constants.CREATED_AT},
			Ccms:        []string{constants.ASSEMBLY_ID, constants.ACCOUNT_ID},
			Hosts:       s.scylla_host,
			Keyspace:    s.scylla_keyspace,
			PksClauses:  map[string]interface{}{constants.EVENT_TYPE: edata.M[constants.EVENT_TYPE], constants.CREATED_AT: s_data.CreatedAt},
			CcmsClauses: map[string]interface{}{constants.ASSEMBLY_ID: edata.M[constants.ASSEMBLY_ID], constants.ACCOUNT_ID: edata.M[constants.ACCOUNT_ID]},
		}	
		if err := ldb.Storedb(ops, s_data); err != nil {
			log.Debugf(err.Error())
			return err
		}
	
	return nil
}

func parseMapToOutputFormat(edata EventData) Events {
   	return Events{
   		EventType:  edata.M[constants.EVENT_TYPE],   
		AccountId:  edata.M[constants.ACCOUNT_ID],
		AssemblyId: edata.M[constants.ASSEMBLY_ID],	
		Data: edata.D,
		CreatedAt: time.Now().String(),
   	}
}




