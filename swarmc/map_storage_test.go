package cluster_test

import (
	"testing"

	"github.com/megamsys/megamd/swarmc"
	 storageTesting "github.com/megamsys/swarmc/storage_testing"
)

func TestMapStorageStorage(t *testing.T) {
	storageTesting.RunTestsForStorage(&cluster.MapStorage{}, t)
}
