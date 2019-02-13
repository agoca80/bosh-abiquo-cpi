package actions

import (
	"github.com/abiquo/ojal/core"
	"github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cloudfoundry/bosh-utils/fileutil"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/bosh-utils/system"
	"github.com/cloudfoundry/bosh-utils/uuid"

	"github.com/cppforlife/bosh-cpi-go/apiv1"

	"github.com/agoca80/bosh-abiquo-cpi/compute"
	"github.com/agoca80/bosh-abiquo-cpi/helpers"
	"github.com/agoca80/bosh-abiquo-cpi/storage"
	"github.com/agoca80/bosh-abiquo-cpi/template"
)

// Factory ...
type Factory struct {
	compressor fileutil.Compressor
	filesystem system.FileSystem
	cmdRunner  system.CmdRunner
	uuidGen    uuid.Generator
	logger     *helpers.Logger
	Options
}

// Options ...
type Options struct {
	DatacenterRepository string
	Endpoint             string
	Username             string
	Password             string
	Agent                apiv1.AgentOptions
}

var _ apiv1.CPIFactory = Factory{}

// CPI ...
type CPI struct {
	Misc
	Stemcells
	VMs
	Disks
	Snapshots
}

var _ apiv1.CPI = CPI{}

// NewFactory ...
func NewFactory(
	filesystem system.FileSystem,
	cmdRunner system.CmdRunner,
	uuidGen uuid.Generator,
	compressor fileutil.Compressor,
	logger logger.Logger,
	options Options,
) Factory {
	err := core.Init(options.Endpoint, core.Basic{
		Username: options.Username,
		Password: options.Password,
	})
	if err != nil {
		panic("could not init the abiquo client")
	}

	return Factory{
		cmdRunner:  cmdRunner,
		compressor: compressor,
		filesystem: filesystem,
		logger:     helpers.NewLogger(logger, "actions.Factory"),
		Options:    options,
		uuidGen:    uuidGen,
	}
}

// New ...
func (f Factory) New(ctx apiv1.CallContext) (apiv1.CPI, error) {
	stemcells := template.NewFactory(
		f.logger.Logger,
		f.uuidGen,
		f.compressor,
		f.filesystem,
		template.Options{
			DatacenterRepository: f.Options.DatacenterRepository,
		},
	)

	disks := storage.NewFactory(
		f.logger.Logger,
		f.uuidGen,
		storage.Options{},
	)

	vms := compute.NewFactory(
		f.logger.Logger,
		f.uuidGen,
		compute.Options{
			Agent: f.Options.Agent,
		},
	)

	return CPI{
		NewMisc(),
		NewStemcells(stemcells, stemcells),
		NewVMs(stemcells, vms, vms),
		NewDisks(disks, disks, vms),
		NewSnapshots(),
	}, nil
}

// Validate ...
func (o Options) Validate() error {
	if o.Username == "" {
		return errors.Error("Must provide non-empty Username")
	}

	if o.Password == "" {
		return errors.Error("Must provide non-empty Password")
	}

	if o.Endpoint == "" {
		return errors.Error("Must provide non-empty Endpoint")
	}

	if o.DatacenterRepository == "" {
		return errors.Error("Must provide non-empty DatacenterRepository href")
	}

	err := o.Agent.Validate()
	if err != nil {
		return errors.WrapError(err, "Validating Agent configuration")
	}

	return nil
}
