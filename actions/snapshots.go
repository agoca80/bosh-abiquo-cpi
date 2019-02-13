package actions

import (
	"github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// Snapshots ...
type Snapshots struct{}

// NewSnapshots ...
func NewSnapshots() Snapshots {
	return Snapshots{}
}

// SnapshotDisk ...
func (s Snapshots) SnapshotDisk(cid apiv1.DiskCID, meta apiv1.DiskMeta) (apiv1.SnapshotCID, error) {
	return apiv1.SnapshotCID{}, errors.Error("unimplemented")
}

// DeleteSnapshot ...
func (s Snapshots) DeleteSnapshot(cid apiv1.SnapshotCID) error {
	return errors.Error("unimplemented")
}
