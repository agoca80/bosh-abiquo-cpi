package compute

import (
	"github.com/abiquo/ojal/core"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/bosh-utils/uuid"
	"github.com/cppforlife/bosh-cpi-go/apiv1"

	"github.com/agoca80/bosh-abiquo-cpi/helpers"
	"github.com/agoca80/bosh-abiquo-cpi/template"
)

// Options ...
type Options struct {
	Agent apiv1.AgentOptions
}

// Factory ...
type Factory struct {
	uuid.Generator
	*helpers.Logger
	Options
}

// NewFactory ...
func NewFactory(
	logger logger.Logger,
	uuidGen uuid.Generator,
	options Options,
) Factory {
	return Factory{
		Generator: uuidGen,
		Options:   options,
		Logger:    helpers.NewLogger(logger, "compute.Factory"),
	}
}

// Create ...
func (f Factory) Create(
	agentID apiv1.AgentID,
	stemcell template.Stemcell,
	props apiv1.VMCloudProps,
	networks apiv1.Networks,
	env apiv1.VMEnv,
) (_ VM, err error) {
	properties, err := newProperties(props)
	if err != nil {
		f.Debug("getting properties: %v", err)
		return
	}

	uuid, err := f.Generate()
	if err != nil {
		return
	}

	vmt, err := helpers.VirtualMachineTemplate(stemcell.ID().AsString())
	if err != nil {
		return
	}

	vapp, err := helpers.VirtualAppliance(properties.VirtualAppliance)
	if err != nil {
		return
	}

	vdc, err := vapp.Rel("virtualdatacenter").Walk()
	if err != nil {
		return
	}

	nets, err := newNetworks(networks)
	if err != nil {
		return
	}

	var hd core.Resource
	if properties.HardDisk > 0 {
		hd, err = helpers.CreateHD(properties.HardDisk, "VIRTIO", vdc)
		if err != nil {
			return
		}
		f.Debug("hd %v", hd.URL())
	}

	f.Debug("creating vm ...")
	vm, err := helpers.CreateVM(uuid, properties.CPUs, properties.Memory, hd, vmt, vapp)
	if err != nil {
		return
	}
	newVM := newVM(vm, f.Logger)
	rollback := func() (VM, error) {
		f.Debug(newVM.Name+" rollback due to %s", err)
		newVM.Delete()
		return nil, err
	}

	f.Debug(newVM.Name + " configuring agent")
	initialAgentEnv := apiv1.NewAgentEnvFactory().ForVM(agentID, newVM.ID(), nets.AsNetworks(), env, f.Options.Agent)
	initialAgentEnv.AttachSystemDisk(apiv1.NewDiskHintFromString("/dev/vda"))
	if properties.HardDisk > 0 {
		initialAgentEnv.AttachEphemeralDisk(apiv1.NewDiskHintFromString("/dev/vdb"))
	}
	err = newVM.bootstrap(initialAgentEnv)
	if err != nil {
		return rollback()
	}

	for _, net := range nets {
		err = newVM.configureNIC(net)
		if err != nil {
			return rollback()
		}
	}

	f.Debug(newVM.Name + " deploying")
	err = newVM.Deploy()
	if err != nil {
		return rollback()
	}

	if synch := newVM.synchronize(); synch != nil {
		f.Debug(newVM.Name+" CreateVM: vm was not synchronized: %v", synch)
	}

	return newVM, nil
}

// Find ...
func (f Factory) Find(cid apiv1.VMCID) (v VM, err error) {
	f.Debug("looking for %s", cid.AsString())
	virtualMachine, err := helpers.VirtualMachine(cid.AsString())
	if err != nil {
		return
	}
	v = newVM(virtualMachine, f.Logger)
	return
}
