// Package db encapsulates connection with Riak.
//
// The function Open dials to Riak and returns a connection (represented by
// the Storage type). It manages an internal pool of connections, and
// reconnects in case of failures. That means that you should not store
// references to the connection, but always call Open.
package db

import (
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/riakpbc"
)

const period time.Duration = 7 * 24 * time.Hour

var (
	conn   = make(map[string]*session) // pool of connections
	mut    sync.RWMutex                // for pool thread safety
	ticker *time.Ticker                // for garbage collection
)

type session struct {
	s    *riakpbc.Client
	used time.Time
}

type RiakDB struct {
	BindAddress []string
	Bucket      string
}

// Storage holds the connection with the bucket name.
type Storage struct {
	coder_client *riakpbc.Client
	bktname      string
}

//a function that returns a new riakdb handler when provided the server bind address and the bucket name.
func NewRiakDB(baddr []string, b string) (*RiakDB, error) {
	return &RiakDB{BindAddress: baddr, Bucket: b}, nil
}

// Conn reads the riak object nd calls Open to get a database connection.
//
// Most megam packages should probably use this function.
func (r *RiakDB) Conn() (*Storage, error) {
	return Open(r.BindAddress, r.Bucket)
}

// Open dials to the Riak database, and return a new connection (represented
// by the type Storage).
//
// addr is a Riak connection URI, and bktname is the name of the bucket.
//
// This function returns a pointer to a Storage, or a non-nil error in case of
// any failure.
func Open(addr []string, bktname string) (storage *Storage, err error) {
	defer func() {
		if r := recover(); r != nil {
			storage, err = open(addr, bktname)
		}
	}()
	mut.RLock()
	if session, ok := conn[strings.Join(addr, "::")]; ok {
		mut.RUnlock()
		if _, err = session.s.Ping(); err == nil {
			mut.Lock()
			session.used = time.Now()
			conn[strings.Join(addr, "::")] = session
			mut.Unlock()
		}
		return open(addr, bktname)
	}
	mut.RUnlock()
	return open(addr, bktname)
}

// Close closes the storage, releasing the connection.
func (s *Storage) Close() {
	log.Debugf(cmd.Colorfy("  > [riak] close", "blue", "", "bold"))
	s.coder_client.Close()
}

// FetchStruct stores a struct  as JSON
//   eg: data := ExampleData{
//        Field1: "ExampleData1",
//        Field2: 1,
//   }
// So the send can pass in 	out := &ExampleData{}
// Apps returns the apps collection from MongoDB.
func (s *Storage) FetchStruct(key string, out interface{}) error {
	if _, err := s.coder_client.FetchStruct(s.bktname, key, out); err != nil {
		return fmt.Errorf("Failed to fetch structure from riak.	--> %s", err.Error())
	}
	//TO-DO:
	//we need to return the fetched json -> to struct interface
	return nil
}

// StoreStruct returns the apps collection from MongoDB.
func (s *Storage) StoreStruct(key string, data interface{}) error {
	if _, err := s.coder_client.StoreStruct(s.bktname, key, data); err != nil {
		return fmt.Errorf("Failed to store a structure in riak --> %s", err.Error())
	}
	return nil
}

func open(addr []string, bucketname string) (*Storage, error) {
	coder := riakpbc.NewCoder("json", riakpbc.JsonMarshaller, riakpbc.JsonUnmarshaller)
	riakCoder := riakpbc.NewClientWithCoder(addr, coder)
	if err := riakCoder.Dial(); err != nil {
		return nil, err
	}
	// Set Client ID
	/*if _, err := riakCoder.SetClientId("coolio"); err != nil {
		log.Fatalf("Setting client ID failed: %v", err)
	}
	*/
	storage := &Storage{coder_client: riakCoder, bktname: bucketname}
	mut.Lock()
	conn[strings.Join(addr, "::")] = &session{s: riakCoder, used: time.Now()}
	mut.Unlock()
	log.Debugf(cmd.Colorfy("  > [riak] open  ", "blue", "", "bold") + fmt.Sprintf("%v success", addr))
	return storage, nil
}

type SomeObject struct {
	Data string
}

// Fetch raw data (int, string, []byte)
func (s *Storage) FetchObject(key string, out *SomeObject) error {
	obj, err := s.coder_client.FetchObject(s.bktname, key)
	if err != nil {
		return fmt.Errorf("Failed to fetch from riak. %s", err.Error())
	}
	out.Data = string(obj.GetContent()[0].GetValue())
	return nil
}

// Store raw data (int, string, []byte)
func (s *Storage) StoreObject(key string, data string) error {
	if _, err := s.coder_client.StoreObject(s.bktname, key, []byte(data)); err != nil {
		return fmt.Errorf("Failed to store in riak. %s", err)
	}
	return nil
}

func (s *Storage) DeleteObject(key string) error {
	if _, err := s.coder_client.DeleteObject(s.bktname, key); err != nil {
		return fmt.Errorf("Failed to delee in riak. %s", err)
	}
	return nil
}

func (s *Storage) FetchObjectByIndex(bucket, index, key, start, end string) ([][]byte, error) {
	number, err := s.coder_client.Index(bucket, index, key, start, end)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch object from riak. %s", err.Error())
	}
	return number.GetKeys(), nil
}

func init() {
	ticker = time.NewTicker(time.Hour)
	go retire(ticker)
}

// retire retires old connections
func retire(t *time.Ticker) {
	for _ = range t.C {
		now := time.Now()
		var old []string
		mut.RLock()
		for k, v := range conn {

			if now.Sub(v.used) >= period {
				old = append(old, k)
			}
		}
		mut.RUnlock()
		mut.Lock()
		for _, c := range old {
			log.Debugf(cmd.Colorfy("  > [riak] stale ", "blue", "", "bold") + fmt.Sprintf("%v ", c))
			conn[c].s.Close()
			delete(conn, c)
		}
		mut.Unlock()
	}
}
