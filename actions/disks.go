package actions

import (
	"github.com/agoca80/bosh-abiquo-cpi/compute"
	"github.com/agoca80/bosh-abiquo-cpi/storage"
	"github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// Disks ...
type Disks struct {
	creator  storage.Creator
	finder   storage.Finder
	vmFinder compute.Finder
}

// NewDisks ...
func NewDisks(creator storage.Creator, finder storage.Finder, vmFinder compute.Finder) Disks {
	return Disks{creator, finder, vmFinder}
}

// AttachDisk ...
func (disks Disks) AttachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) (err error) {
	_, err = disks.AttachDiskV2(vmCID, diskCID)
	return
}

// AttachDiskV2 ...
func (disks Disks) AttachDiskV2(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) (apiv1.DiskHint, error) {
	fail := func(msg string, err error) (apiv1.DiskHint, error) {
		return apiv1.DiskHint{}, errors.WrapErrorf(err, msg)
	}

	vm, err := disks.vmFinder.Find(vmCID)
	if err != nil {
		return fail("finding VM", err)
	}

	disk, err := disks.finder.Find(diskCID)
	if err != nil {
		return fail("finding disk", err)
	}

	hint, err := vm.AttachDisk(disk)
	if err != nil {
		return fail("attaching disk to VM", err)
	}

	return hint, nil
}

// CreateDisk ...
func (disks Disks) CreateDisk(size int, props apiv1.DiskCloudProps, vmCID *apiv1.VMCID) (id apiv1.DiskCID, err error) {
	disk, err := disks.creator.CreateDisk(size, props)
	if err != nil {
		return
	}

	id = disk.ID()
	return
}

// DeleteDisk ...
func (disks Disks) DeleteDisk(cid apiv1.DiskCID) (err error) {
	disk, err := disks.finder.Find(cid)
	if err != nil {
		return
	}

	return disk.Delete()
}

// DetachDisk ...
func (disks Disks) DetachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) (err error) {
	vm, err := disks.vmFinder.Find(vmCID)
	if err != nil {
		return
	}

	disk, err := disks.finder.Find(diskCID)
	if err != nil {
		return
	}

	err = vm.DetachDisk(disk)
	if err != nil {
		return
	}

	return
}

// GetDisks ...
func (vms VMs) GetDisks(cid apiv1.VMCID) (disks []apiv1.DiskCID, err error) {
	vm, err := vms.finder.Find(cid)
	if err != nil {
		return
	}

	disks, err = vm.DiskIDs()
	return
}

// HasDisk ...
func (disks Disks) HasDisk(cid apiv1.DiskCID) (ok bool, err error) {
	_, err = disks.finder.Find(cid)
	if err != nil {
		return
	}

	ok = true
	return
}

// ResizeDisk ...
func (disks Disks) ResizeDisk(cid apiv1.DiskCID, size int) error {
	return errors.Error("unimplemented")
}

// SetDiskMetadata ...
func (disks Disks) SetDiskMetadata(cid apiv1.DiskCID, metadata apiv1.DiskMeta) (err error) {
	disk, err := disks.finder.Find(cid)
	if err != nil {
		return
	}

	err = disk.SetMetadata(metadata)
	return
}
