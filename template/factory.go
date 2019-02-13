package template

import (
	"path/filepath"

	"github.com/abiquo/ojal/abiquo"
	"github.com/cppforlife/bosh-cpi-go/apiv1"

	"github.com/cloudfoundry/bosh-utils/fileutil"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/bosh-utils/system"
	"github.com/cloudfoundry/bosh-utils/uuid"

	"github.com/agoca80/bosh-abiquo-cpi/helpers"
)

// Factory ...
type Factory struct {
	Options
	*helpers.Logger

	compressor fileutil.Compressor
	filesystem system.FileSystem
	uuidGen    uuid.Generator
}

// Options ...
type Options struct {
	DatacenterRepository string
}

// NewFactory ...
func NewFactory(
	logger logger.Logger,
	uuidGen uuid.Generator,
	compressor fileutil.Compressor,
	filesystem system.FileSystem,
	options Options,
) Factory {
	return Factory{
		compressor: compressor,
		filesystem: filesystem,
		uuidGen:    uuidGen,
		Options:    options,
		Logger:     helpers.NewLogger(logger, "stemcell.Factory"),
	}
}

// ImportFromPath ...
func (f Factory) ImportFromPath(path string) (stemcell Stemcell, err error) {
	f.Debug("generating template uuid")
	uuid, err := f.uuidGen.Generate()
	if err != nil {
		return
	}

	f.Debug("getting datacenter repository")
	datacenterRepository, err := helpers.DatacenterRepository(f.DatacenterRepository)
	if err != nil {
		return
	}

	f.Debug("creating temporary directory")
	tempDir, err := f.filesystem.TempDir("upload-stemcell")
	if err != nil {
		return
	}
	defer f.filesystem.RemoveAll(tempDir)

	f.Debug("decompressing temporary directory")
	err = f.compressor.DecompressFileToDir(path, tempDir, fileutil.CompressorOptions{})
	if err != nil {
		return
	}

	f.Debug("uploading root disk")
	filename := filepath.Join(tempDir, "root.img")
	vmt, err := datacenterRepository.UploadTemplate(filename, abiquo.TemplateDefinition{
		CategoryName: "Others",
		Name:         uuid,
		Disks: []abiquo.DiskDefinition{
			abiquo.DiskDefinition{
				Bootable:           true,
				DiskControllerType: "VIRTIO",
				DiskFileFormat:     "QCOW2_SPARSE",
				DiskFilePath:       filename,
				DiskFileSize:       500000000,
				RequiredHDinMB:     3072,
				Sequence:           0,
			},
		},
		RequiredCPU:     1,
		RequiredRAMInMB: 512,
	})
	if err != nil {
		return
	}

	f.Debug("updating template %s", vmt.URL())
	vmt.IconURL = "https://www.cloudfoundry.org/wp-content/uploads/2017/10/CFF-Symbol-BOSH-Full-Color.png"
	vmt.Name = uuid
	vmt.Description = "bosh stemcell"
	vmt.EthernetDriverType = "VIRTIO"
	vmt.GuestSetup = "CLOUD_INIT"
	vmt.GenerateGuestInitialPassword = true
	vmt.LoginUser = "abiquo"
	err = vmt.Update(vmt)
	if err != nil {
		return
	}

	stemcell = newTemplate(vmt, f.Logger)
	return
}

// Find ...
func (f Factory) Find(cid apiv1.StemcellCID) (stemcell Stemcell, err error) {
	f.Debug("looking for %s", cid.AsString())
	vmt, err := helpers.VirtualMachineTemplate(cid.AsString())
	if err != nil {
		return
	}
	stemcell = newTemplate(vmt, f.Logger)
	return
}
