package events

// Interface for event  operation handlers.
type EventHandler interface {
	// Returns the ContainerReference
	//VMReference() (info.ContainerReference, error) //*EventHoldingReference - can be an interface which get casted ?

	// Registers a channel to listen for events affecting vm(eventholder) (recursively).
	Watch(events chan Event) error

	// Stops watching for changes.
	StopWatching() error

	// Returns whether the vm(type of eventholder) still exists.
	Exists() bool

	Cleanup()

	// Start starts any necessary background goroutines - must be cleaned up in Cleanup().
	// It is expected that most implementations will be a no-op.
	Start()
}
