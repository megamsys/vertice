package swarmc

import (
	"testing"

	 storageTesting "github.com/megamsys/swarmc/storage_testing"
)

func TestMapStorageStorage(t *testing.T) {
	storageTesting.RunTestsForStorage(&cluster.MapStorage{}, t)
}
