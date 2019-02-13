package actions

import (
	"github.com/agoca80/bosh-abiquo-cpi/compute"
	"github.com/agoca80/bosh-abiquo-cpi/template"
	"github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// VMs ...
type VMs struct {
	stemcells template.Finder
	creator   compute.Creator
	finder    compute.Finder
}

// NewVMs ...
func NewVMs(stemcells template.Finder, creator compute.Creator, finder compute.Finder) VMs {
	return VMs{stemcells, creator, finder}
}

// CalculateVMCloudProperties ...
func (vms VMs) CalculateVMCloudProperties(res apiv1.VMResources) (apiv1.VMCloudProps, error) {
	return apiv1.NewVMCloudPropsFromMap(map[string]interface{}{
		"memory":         res.RAM,
		"cpus":           res.CPU,
		"ephemeral_disk": res.EphemeralDiskSize,
	}), nil
}

// CreateVM ...
func (vms VMs) CreateVM(
	agentID apiv1.AgentID,
	stemcellCID apiv1.StemcellCID,
	props apiv1.VMCloudProps,
	networks apiv1.Networks,
	diskCIDs []apiv1.DiskCID,
	env apiv1.VMEnv,
) (vmcid apiv1.VMCID, err error) {
	cid, _, err := vms.CreateVMV2(agentID, stemcellCID, props, networks, diskCIDs, env)
	if err != nil {
		return cid, errors.WrapErrorf(err, "Creating vm '%s'", err.Error())
	}

	return cid, err
}

// CreateVMV2 ...
func (vms VMs) CreateVMV2(
	agentID apiv1.AgentID,
	stemcellCID apiv1.StemcellCID,
	props apiv1.VMCloudProps,
	networks apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID,
	env apiv1.VMEnv,
) (cid apiv1.VMCID, nets apiv1.Networks, err error) {
	stemcell, err := vms.stemcells.Find(stemcellCID)
	if err != nil {
		return apiv1.VMCID{}, networks, errors.WrapErrorf(err, "Finding stemcell '%s'", stemcellCID)
	}

	vm, err := vms.creator.Create(agentID, stemcell, props, networks, env)
	if err != nil {
		return apiv1.VMCID{}, networks, errors.WrapErrorf(err, "Creating VM with agent ID '%s'", agentID)
	}

	return vm.ID(), networks, nil
}

// DeleteVM ...
func (vms VMs) DeleteVM(cid apiv1.VMCID) error {
	vm, err := vms.finder.Find(cid)
	if err != nil {
		return errors.WrapErrorf(err, "Finding vm '%s'", cid)
	}

	err = vm.Delete()
	if err != nil {
		return errors.WrapErrorf(err, "Deleting vm '%s'", cid)
	}

	return nil
}

// HasVM ...
func (vms VMs) HasVM(cid apiv1.VMCID) (bool, error) {
	_, err := vms.finder.Find(cid)
	if err != nil {
		err = errors.WrapErrorf(err, "Finding VM %s", cid.AsString())
	}

	return err == nil, err
}

// RebootVM ...
func (vms VMs) RebootVM(cid apiv1.VMCID) error {
	vm, err := vms.finder.Find(cid)
	if err != nil {
		return errors.WrapErrorf(err, "Finding VM %s", cid.AsString())
	}

	return vm.Reboot()
}

// SetVMMetadata ...
// This is currently a NOP
func (vms VMs) SetVMMetadata(cid apiv1.VMCID, metadata apiv1.VMMeta) (err error) {
	vm, err := vms.finder.Find(cid)
	if err != nil {
		return errors.WrapErrorf(err, "Finding VM %s", cid.AsString())
	}

	err = vm.SetMetadata(metadata)
	return
}
