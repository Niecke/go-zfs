package zfs

import (
	"testing"

	zfs "github.com/Niecke/go-zfs"
)

func TestListZpools(t *testing.T) {
	pools, err := zfs.ListZpools()
	if err != nil {
		t.Fatalf(`%v`, err)
	}
	for _, pool := range pools {
		t.Logf("%+v\n", pool)
	}
	if len(pools) != 1 {
		t.Errorf("%d pools in list should be 1", len(pools))
	}
}
