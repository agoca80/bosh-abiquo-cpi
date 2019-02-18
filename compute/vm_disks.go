package compute

import (
	"github.com/abiquo/ojal/core"
	"github.com/agoca80/bosh-abiquo-cpi/storage"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// DiskIDs ...
func (v *vm) DiskIDs() ([]apiv1.DiskCID, error) {
	var persistentIDs []apiv1.DiskCID
	for _, disk := range v.Disks() {
		if disk.IsMedia("volume") {
			persistentIDs = append(persistentIDs, apiv1.NewDiskCID(disk.URL()))
		}
	}
	return persistentIDs, nil
}

// DetachDisk ...
func (v *vm) DetachDisk(disk storage.Disk) (err error) {
	return v.reconfigure(func() error {
		v.Debug(v.UUID+" detaching %s", disk.ID())
		v.VirtualMachine.DetachDisk(disk.Disk())
		return nil
	})
}

// AttachDisk ...
func (v *vm) AttachDisk(disk storage.Disk) (hint apiv1.DiskHint, err error) {
	err = v.reconfigure(func() error {
		v.Debug(v.UUID+" attaching %s", disk.ID().AsString())
		v.VirtualMachine.AttachDisk(disk.Disk())
		return nil
	})
	if err != nil {
		return
	}

	disks := v.Disks().Filter(func(l *core.Link) bool {
		return l.DiskControllerType == "VIRTIO"
	})
	hint = apiv1.NewDiskHintFromMap(map[string]interface{}{
		"path": "/dev/vd" + string('a'+len(disks)-1),
	})
	return
}
