package storage

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/agoca80/bosh-abiquo-cpi/helpers"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type disk struct {
	*abiquo.Volume
	*helpers.Logger
}

// ID returns the disk href
func (d *disk) ID() apiv1.DiskCID {
	d.Debug("volume %s", d.URL())
	return apiv1.NewDiskCID((d.URL()))
}

// Exists returns if the disk still exists in the Abiquo backend
func (d *disk) Exists() (exists bool, err error) {
	return d.Volume.Exists()
}

// Delete removes the disk's Hard Disk resource from Abiquo
func (d *disk) Delete() error {
	return d.Volume.Remove()
}

// SetMetadata ...
// This is a NOP as Abiquo does not provide any way to provide harddisks metadata
// Perhaps name/description field for volumes?
func (d *disk) SetMetadata(metadata apiv1.DiskMeta) (err error) {
	d.Debug("SetMetadata: metadata %+v", metadata)
	marshalled, err := metadata.MarshalJSON()
	if err != nil {
		return
	}
	d.Debug("SetMetadata: %v", string(marshalled))
	// This may be longer that 255 characters...
	// d.Description = base64.StdEncoding.EncodeToString(marshalled)
	// err = d.Volume.Update()
	return
}

// Disk ...
func (d *disk) Disk() abiquo.Disk {
	return d.Volume
}

// ResizeDisk ...
// func (d *disk) ResizeDisk(cid apiv1.DiskCID, size int) (err error) {
// 	d.Debug("finding volume")
// 	disk, err := d.Find(cid)
// 	if err != nil {
// 		return
// 	}

// 	d.Debug("resizing volume")
// 	disk.SizeInMB = size
// 	return disk.Update()
// }

func newDisk(volume *abiquo.Volume, logger *helpers.Logger) (d *disk) {
	d = &disk{volume, logger}
	return
}
