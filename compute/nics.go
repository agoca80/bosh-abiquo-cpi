package compute

import (
	"fmt"

	"github.com/abiquo/ojal/core"
)

func (v *vm) adquireIP(network *Network) (err error) {
	v.Debug(v.Name+" adquireIP: getting IP from %q", network.Network.Name)
	link, err := getLink[network.API.Type()](network)
	if err != nil {
		v.Debug(v.Name+" adquireIP: error %v", err)
		return
	}

	return v.change(func() error {
		v.Debug(v.UUID+" adquireIP: adquiring ip %v", link.Title)
		v.AttachNIC(link)
		return nil
	})
}

func (v *vm) configureNIC(network *Network) (err error) {
	v.Debug(v.UUID+" configureNIC: configuring %s IP in %q", network.API.Type(), network.Network.Name)

	canRetry := func(err error) (can bool) {
		e, ok := err.(core.Error)
		if ok && len(e.Collection) == 1 {
			can = e.Collection[0].Code == "VM-24" || e.Collection[0].Code == "DB-0"
		}
		return
	}

	for try := 0; try < 6; try++ {
		err = v.adquireIP(network)
		if err == nil || !canRetry(err) {
			return
		}

		// Refresh virtual machine
		err = v.Read(v.VirtualMachine)
		if err != nil {
			return
		}
	}

	return fmt.Errorf(v.UUID+" configureNIC: couldnt adquire IP in %q", network.Network.Name)
}
