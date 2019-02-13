package storage

import (
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/bosh-utils/uuid"
	"github.com/cppforlife/bosh-cpi-go/apiv1"

	"github.com/agoca80/bosh-abiquo-cpi/helpers"
)

// Factory ...
type Factory struct {
	uuid.Generator
	*helpers.Logger
	Options
}

// Options ...
type Options struct{}

// NewFactory ...
func NewFactory(
	logger logger.Logger,
	uuidGen uuid.Generator,
	options Options,
) (factory Factory) {
	return Factory{
		Logger:    helpers.NewLogger(logger, "storage.Factory"),
		Options:   options,
		Generator: uuidGen,
	}
}

// Find ...
func (f Factory) Find(cid apiv1.DiskCID) (d Disk, err error) {
	f.Debug("looking for %s", cid.AsString())
	volume, err := helpers.Volume(cid.AsString())
	if err != nil {
		return
	}
	return newDisk(volume, f.Logger), nil
}

// CreateDisk ...
func (f Factory) CreateDisk(size int, props apiv1.DiskCloudProps) (disk Disk, err error) {
	f.Debug("getting properties")
	properties, err := f.newProperties(props)
	if err != nil {
		return
	}

	f.Debug("generating uuid")
	uuid, err := f.Generate()
	if err != nil {
		return
	}

	f.Debug("getting vdc")
	vdc, err := helpers.VirtualDatacenter(properties.VDC)
	if err != nil {
		return
	}

	f.Debug("getting tier")
	tier, err := helpers.Tier(properties.Tier)
	if err != nil {
		return
	}

	f.Debug("creating volume")
	volume, err := helpers.CreateVolume(size, uuid, "VIRTIO", tier, vdc)
	if err != nil {
		return
	}
	return newDisk(volume, f.Logger), nil
}
