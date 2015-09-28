package repository

func init() {
	Register("nop", nopManager{})
}

type nopManager struct{}

func (nopManager) CreateHook(r Repository) (string, error) {
	return "", nil
}

func (nopManager) RemoveHook(r Repository) error {
	return nil
}
