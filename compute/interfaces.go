package compute

import (
	"github.com/agoca80/bosh-abiquo-cpi/storage"
	"github.com/agoca80/bosh-abiquo-cpi/template"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// Creator ...
type Creator interface {
	Create(
		apiv1.AgentID,
		template.Stemcell,
		apiv1.VMCloudProps,
		apiv1.Networks,
		apiv1.VMEnv,
	) (VM, error)
}

var _ Creator = Factory{}

// Finder ...
type Finder interface {
	Find(apiv1.VMCID) (VM, error)
}

var _ Finder = Factory{}

// VM ...
type VM interface {
	ID() apiv1.VMCID
	SetMetadata(apiv1.VMMeta) error

	Reboot() error
	Exists() (bool, error)
	Delete() error

	DiskIDs() ([]apiv1.DiskCID, error)
	AttachDisk(storage.Disk) (apiv1.DiskHint, error)
	DetachDisk(storage.Disk) error
}

var _ VM = &vm{}
