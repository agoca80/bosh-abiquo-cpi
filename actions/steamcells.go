package actions

import (
	"github.com/agoca80/bosh-abiquo-cpi/helpers"
	"github.com/agoca80/bosh-abiquo-cpi/template"
	"github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// Stemcells ...
type Stemcells struct {
	importer template.Importer
	finder   template.Finder
}

// NewStemcells ...
func NewStemcells(importer template.Importer, finder template.Finder) Stemcells {
	return Stemcells{importer, finder}
}

// CreateStemcell ...
func (s Stemcells) CreateStemcell(
	imagePath string,
	props apiv1.StemcellCloudProps,
) (cid apiv1.StemcellCID, err error) {
	properties := template.Properties{}
	err = props.As(&properties)
	if err != nil {
		return
	}

	// Ugly crap
	var stemcell template.Stemcell
	if properties.Href != "" {
		helpers.Debug("stemcell href is %s", properties.Href)
		stemcell, err = s.finder.Find(apiv1.NewStemcellCID(properties.Href))
		if err != nil {
			return
		}

		return stemcell.ID(), nil
	}

	stemcell, err = s.importer.ImportFromPath(imagePath)
	if err != nil {
		return apiv1.StemcellCID{}, errors.WrapErrorf(err, "Importing stemcell from '%s'", imagePath)
	}

	return stemcell.ID(), nil
}

// DeleteStemcell ...
func (s Stemcells) DeleteStemcell(cid apiv1.StemcellCID) error {
	stemcell, err := s.finder.Find(cid)
	if err != nil {
		return errors.WrapErrorf(err, "finding stemcell %s", cid.AsString())
	}

	return stemcell.Delete()
}
