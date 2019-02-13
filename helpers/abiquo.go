package helpers

import (
	"github.com/abiquo/ojal/abiquo"
	"github.com/abiquo/ojal/core"
)

// CreateHD ...
func CreateHD(size int, ctrlType string, vdc core.Resource) (hd *abiquo.HardDisk, err error) {
	hd = &abiquo.HardDisk{
		DiskControllerType: ctrlType,
		SizeInMb:           size,
	}
	err = vdc.Rel("disks").SetType("harddisk").Create(hd)
	return
}

// CreateVM ...
func CreateVM(uuid string, cpus, ram int, hd, vmt, vapp core.Resource) (vm *abiquo.VirtualMachine, err error) {
	vm = &abiquo.VirtualMachine{
		Label:       uuid,
		CPU:         cpus,
		RAM:         ram,
		Password:    "12qwaszx",
		VdrpEnabled: true,
		DTO: core.NewDTO(
			hd.Link().SetRel("disk1"),
			vmt.Link().SetRel("virtualmachinetemplate"),
		),
	}
	err = vapp.Rel("virtualmachines").SetType("virtualmachine").Create(vm)
	return
}

// CreateVolume ...
func CreateVolume(size int, uuid, ctrlType string, tier, vdc core.Resource) (volume *abiquo.Volume, err error) {
	volume = &abiquo.Volume{
		DiskControllerType: ctrlType,
		Name:               uuid,
		SizeInMB:           size,
		DTO: core.NewDTO(
			tier.Link().SetRel("tier"),
		),
	}
	err = vdc.Rel("volumes").SetType("volume").Create(volume)
	return
}

// VirtualMachineTemplate ...
func VirtualMachineTemplate(href string) (resource *abiquo.VirtualMachineTemplate, err error) {
	resource = new(abiquo.VirtualMachineTemplate)
	err = core.NewLink(href).SetType("virtualmachinetemplate").Read(resource)
	return
}

// DatacenterRepository ...
func DatacenterRepository(href string) (resource *abiquo.DatacenterRepository, err error) {
	resource = new(abiquo.DatacenterRepository)
	err = core.NewLink(href).SetType("datacenterrepository").Read(resource)
	return
}

// Network ...
func Network(href string) (resource *abiquo.Network, err error) {
	resource = new(abiquo.Network)
	err = core.NewLink(href).SetType("tier").Read(resource)
	return
}

// Tier ...
func Tier(href string) (resource *abiquo.Tier, err error) {
	resource = new(abiquo.Tier)
	err = core.NewLink(href).SetType("tier").Read(resource)
	return
}

// VirtualAppliance ...
func VirtualAppliance(href string) (resource *abiquo.VirtualAppliance, err error) {
	resource = new(abiquo.VirtualAppliance)
	err = core.NewLink(href).SetType("virtualappliance").Read(resource)
	return
}

// VirtualDatacenter ...
func VirtualDatacenter(href string) (resource *abiquo.VirtualDatacenter, err error) {
	resource = new(abiquo.VirtualDatacenter)
	err = core.NewLink(href).SetType("virtualdatacenter").Read(resource)
	return
}

// VirtualMachine ...
func VirtualMachine(href string) (resource *abiquo.VirtualMachine, err error) {
	resource = new(abiquo.VirtualMachine)
	err = core.NewLink(href).SetType("virtualmachine").Read(resource)
	return
}

// Volume ...
func Volume(href string) (resource *abiquo.Volume, err error) {
	resource = new(abiquo.Volume)
	err = core.NewLink(href).SetType("volume").Read(resource)
	return
}
