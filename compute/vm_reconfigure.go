package compute

import (
	"github.com/abiquo/ojal/abiquo"
)

func (v *vm) requiresStop() bool {
	// The only reconfigure operation triggered when a VM is running
	// is the attach/dettach disk operation. Hence, we only need to check
	// for EnableDiskHotReconfigure

	vmt := new(abiquo.VirtualMachineTemplate)
	err := v.Rel("virtualmachinetemplate").Read(vmt)
	if err != nil {
		v.Debug(v.Name + " vmt not found. Assuming hot reconfigure support is false")
		return v.State == "ON"
	}

	return v.State == "ON" && !vmt.EnableDisksHotReconfigure
}

func (v *vm) change(fn func() error) (err error) {
	v.Debug(v.Name + " performing changes")
	err = fn()
	if err != nil {
		return
	}

	v.Debug(v.Name + " reconfiguring vm")
	return v.Reconfigure()
}

func (v *vm) hotReconfigure(fn func() error) (err error) {
	v.Debug(v.Name + " performing hotReconfigure")
	v.change(fn)

	if synch := v.synchronize(); synch != nil {
		v.Debug(v.Name+" hotReconfigure: vm was not synchronized: %v", synch)
	}
	return
}

// Reconfigure ...
// TODO: Take into account hot-reconfigure cappable azs...
func (v *vm) reconfigure(fn func() error) (err error) {
	switch v.requiresStop() {
	case true:
		err = v.coldReconfigure(fn)
	case false:
		err = v.hotReconfigure(fn)
	}
	return
}

func (v *vm) coldReconfigure(fn func() error) (err error) {
	v.Debug(v.Name + " performing coldReconfigure")
	err = v.Shutdown()
	if err != nil {
		return
	}

	err = v.change(fn)
	if err != nil {
		return
	}

	if synch := v.synchronize(); synch != nil {
		v.Debug(v.Name+" coldReconfigure: vm was not synchronized: %v", synch)
	}

	return v.On()
}
