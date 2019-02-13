package compute

import (
	"encoding/json"

	"github.com/abiquo/ojal/abiquo"
	"github.com/agoca80/bosh-abiquo-cpi/helpers"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// VMImpl represents an Abiquo VM DTO in the cpi backend
type vm struct {
	*helpers.Logger
	*abiquo.VirtualMachine
}

func newVM(virtualMachine *abiquo.VirtualMachine, logger *helpers.Logger) (v *vm) {
	return &vm{
		Logger:         logger,
		VirtualMachine: virtualMachine,
	}
}

// ID returns the VM cloud ID
func (v *vm) ID() apiv1.VMCID {
	return apiv1.NewVMCID(v.URL())
}

// Delete ...
func (v *vm) Delete() (err error) {
	err = v.VirtualMachine.Delete()
	if err != nil {
		v.Debug("while deleting %v : %v", v.Name, err)
	}
	return
}

// SetMetadata ...
func (v *vm) SetMetadata(metadata apiv1.VMMeta) (err error) {
	v.Debug(v.UUID+" Setmetadata: %+v", metadata)
	marshalled, err := json.Marshal(metadata)
	if err != nil {
		return
	}
	v.Debug("SetMetadata: %v", string(marshalled))

	err = json.Unmarshal(marshalled, &v.Variables)
	if err != nil {
		return
	}

	return v.Reconfigure()
}

// Exists ...
func (v *vm) Exists() (bool, error) {
	return v.VirtualMachine.Exists()
}

// Reboot ...
func (v *vm) Reboot() error {
	return v.VirtualMachine.Reboot()
}
