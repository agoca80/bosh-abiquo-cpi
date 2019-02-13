package template

import (
	apiv1 "github.com/cppforlife/bosh-cpi-go/apiv1"
)

// Importer ...
type Importer interface {
	ImportFromPath(string) (Stemcell, error)
}

var _ Importer = Factory{}

// Finder ...
type Finder interface {
	Find(apiv1.StemcellCID) (Stemcell, error)
}

var _ Finder = Factory{}

// Stemcell ...
type Stemcell interface {
	ID() apiv1.StemcellCID
	Exists() (bool, error)
	Delete() error
}

var _ Stemcell = &stemcell{}
