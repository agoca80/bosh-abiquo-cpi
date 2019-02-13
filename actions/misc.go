package actions

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// Misc ...
type Misc struct{}

// NewMisc ...
func NewMisc() (m Misc) {
	return m
}

// Info ...
func (m Misc) Info() (apiv1.Info, error) {
	return apiv1.Info{
		APIVersion:      2,
		StemcellFormats: []string{"abiquo-qcow2"},
	}, nil
}
