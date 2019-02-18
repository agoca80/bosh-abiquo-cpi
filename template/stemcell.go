package template

import (
	"github.com/agoca80/bosh-abiquo-cpi/helpers"
	apiv1 "github.com/cppforlife/bosh-cpi-go/apiv1"

	"github.com/abiquo/ojal/abiquo"
)

// Resource represents an stemcell's Abiquo VirtualImage
type stemcell struct {
	*abiquo.VirtualMachineTemplate
	*helpers.Logger
}

// NewResource returns a new Abiquo CPI
func newTemplate(vmt *abiquo.VirtualMachineTemplate, logger *helpers.Logger) (s *stemcell) {
	s = &stemcell{
		VirtualMachineTemplate: vmt,
		Logger:                 logger,
	}
	s.Debug("vmt %s", s.URL())
	return
}

// ID returns the stemcell cloud id
func (s *stemcell) ID() apiv1.StemcellCID {
	return apiv1.NewStemcellCID(s.URL())
}

// Exists returns if the disk still exists in the Abiquo backend
func (s *stemcell) Exists() (exists bool, err error) {
	return s.Exists()
}

// Delete removes the stemcell's Abiquo Virtual Image from Abiquo
func (s *stemcell) Delete() error {
	// return core.Remove(r.resource)
	// XXX
	helpers.Debug("NOP: do not delete stemcell during development")
	return nil
}
