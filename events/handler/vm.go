package handler

// VMEventType indicates an addition, deletion or status event.
type VMEventType int

const (
	VMAdd VMEventType = iota
	VMDelete
	VMStatus
)

// VMEvent represents events from a vm (in this case VM is the eventholder)
type VMEvent struct {
	// The type of event that occurred.
	EventType VMEventType

	// The full name of the vm where the event occurred.
	Name string

	//Watcher channel
}

// Interface for container operation handlers.
type VMEventHandler interface {
	// Returns the ContainerReference
	//VMReference() (info.ContainerReference, error) //*EventHoldingReference - can be an interface which get casted ?

	// Registers a channel to listen for events affecting vm(eventholder) (recursively).
	Watch(events chan VMEvent) error

	// Stops watching for changes.
	StopWatching() error

	// Returns whether the vm(type of eventholder) still exists.
	Exists() bool

	Cleanup()

	// Start starts any necessary background goroutines - must be cleaned up in Cleanup().
	// It is expected that most implementations will be a no-op.
	Start()
}

/*just wait on the channel// Watches for new containers started in the system. Runs forever unless there is a setup error.
func (self *manager) Watch(quit chan error) error {
	eventsChannel := make(chan container.VMEvent, 16)
	err := root.handler.Watch(eventsChannel)
	if err != nil {
		return err
	}

	// Listen to events from the container handler.
	go func() {
		for {
			select {
			case event := <-eventsChannel:
				switch {
				case event.EventType == VMAdd:
					err = self.create(event.Name)
				case event.EventType == VMDelete:
					err = self.destroy(event.Name)
				case event.EventType == VMStatus:
					err = self.status(event.Name)
				}
				if err != nil {
					log.Warningf("Failed to process watch event: %v", err)
				}
			case <-quit:
				// Stop processing events if asked to quit.
				err := root.handler.StopWatchingSubcontainers()
				quit <- err
				if err == nil {
					log.Infof("Exiting thread watching subcontainers")
					return
				}
			}
		}
 	}()
  	return nil
}
*/
