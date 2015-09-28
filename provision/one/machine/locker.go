package machine

type BoxLocker interface {
	Lock(appName string) bool
	Unlock(appName string)
}
