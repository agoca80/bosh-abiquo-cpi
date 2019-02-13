package compute

import (
	"errors"

	"github.com/agoca80/bosh-abiquo-cpi/helpers"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// Properties ...
type Properties struct {
	CPUs             int
	HardDisk         int
	Memory           int
	VirtualAppliance string `json:"virtualappliance"`
}

func newProperties(props apiv1.VMCloudProps) (properties *Properties, err error) {
	p := &Properties{}
	hash := make(map[string]interface{})
	if props != nil {
		props.As(&hash)
		props.As(&p)
	}
	helpers.Msg("properties : %+v", hash)

	err = p.validate()
	if err != nil {
		return
	}

	properties = p
	return
}

// TODO. Clean this gruesome error chaining...
func (p *Properties) validate() (err error) {
	msg := ""

	if p.CPUs == 0 {
		msg += "- CPU count cant be 0\n"
	}

	if p.Memory == 0 {
		msg += "- VM RAM count cant be 0\n"
	}

	if p.HardDisk == 0 {
		msg += "- Ephemeral hard disk size cant be 0\n"
	}

	if p.VirtualAppliance == "" {
		msg += "- missing VM virtualappliance property\n"
	}

	if msg != "" {
		msg = "compute.Properties.validate:\n" + msg
		err = errors.New(msg)
	}

	return
}
