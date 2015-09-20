package api

import (
	"sync"

	"github.com/megamsys/megamd/provision"
)

type logStreamTracker struct {
	sync.Mutex
	conn map[*provision.LogListener]struct{}
}

func (t *logStreamTracker) add(l *provision.LogListener) {
	t.Lock()
	defer t.Unlock()
	if t.conn == nil {
		t.conn = make(map[*provision.LogListener]struct{})
	}
	t.conn[l] = struct{}{}
}

func (t *logStreamTracker) remove(l *provision.LogListener) {
	t.Lock()
	defer t.Unlock()
	if t.conn == nil {
		t.conn = make(map[*provision.LogListener]struct{})
	}
	delete(t.conn, l)
}

func (t *logStreamTracker) String() string {
	return "log pub/sub connections"
}

func (t *logStreamTracker) Shutdown() {
	t.Lock()
	defer t.Unlock()
	for l := range t.conn {
		l.Close()
	}
}

var LogTracker logStreamTracker
