package storage

import (
	"errors"

	"github.com/abiquo/ojal/abiquo"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// Properties ...
type Properties struct {
	VDC  string `json:"virtualdatacenter"`
	Tier string `json:"tier"`

	// Dependencies
	size int
	tier *abiquo.Tier
	vdc  *abiquo.VirtualDatacenter
	uuid string
}

func (f Factory) newProperties(props apiv1.DiskCloudProps) (properties *Properties, err error) {
	p := &Properties{}
	hash := make(map[string]interface{})
	if props != nil {
		props.As(&hash)
		props.As(&p)
	}
	f.Debug("properties hash   %+v", hash)
	f.Debug("properties struct %+v", properties)

	err = p.validate()
	if err != nil {
		return
	}

	properties = p
	return
}

// Validate ...
func (p *Properties) validate() (err error) {
	msg := ""

	if p.VDC == "" {
		msg += "- missing virtualdatacenter href\n"
	}

	if p.Tier == "" {
		msg += "- missing tier href\n"
	}

	if msg != "" {
		msg = "storage.Properties.validate:\n" + msg
		err = errors.New(msg)
	}

	return
}
