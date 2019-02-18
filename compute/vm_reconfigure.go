package compute

import (
	"github.com/abiquo/ojal/abiquo"
)

func (v *vm) requiresStop() bool {
	vmt := new(abiquo.VirtualMachineTemplate)
	err := v.Rel("virtualmachinetemplate").Read(vmt)
	if err != nil {
		v.Debug(v.UUID + " vmt not found. Assuming hot reconfigure support is false")
		return v.State == "ON"
	}

	return v.State == "ON" && !vmt.EnableDisksHotReconfigure
}

func (v *vm) change(fn func() error) (err error) {
	v.Debug(v.UUID + " reconfiguring vm")
	err = fn()
	if err != nil {
		return
	}

	err = v.Reconfigure()
	if err != nil {
		v.Debug(v.UUID+" during vm reconfiguration: %v", err)
		return
	}

	if v.State == "NOT_ALLOCATED" {
		return
	}

	return v.Synchronize()
}

func (v *vm) reconfigure(fn func() error) (err error) {
	requiresStop := v.requiresStop()

	if requiresStop {
		err = v.Shutdown()
		if err != nil {
			return
		}
	}

	err = v.change(fn)
	if err != nil {
		return
	}

	if requiresStop {
		err = v.On()
	}

	return
}
