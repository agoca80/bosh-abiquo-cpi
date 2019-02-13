package storage

import (
	"github.com/abiquo/ojal/abiquo"
	apiv1 "github.com/cppforlife/bosh-cpi-go/apiv1"
)

// Creator ...
type Creator interface {
	CreateDisk(int, apiv1.DiskCloudProps) (Disk, error)
}

var _ Creator = Factory{}

// Finder ...
type Finder interface {
	Find(apiv1.DiskCID) (Disk, error)
}

var _ Finder = Factory{}

// Disk ...
type Disk interface {
	ID() apiv1.DiskCID
	Disk() abiquo.Disk
	SetMetadata(apiv1.DiskMeta) error

	// ResizeDisk(apiv1.DiskCID, int) error
	Exists() (bool, error)
	Delete() error
}

var _ Disk = &disk{}
