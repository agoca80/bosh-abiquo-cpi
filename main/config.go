package main

import (
	"encoding/json"

	berr "github.com/cloudfoundry/bosh-utils/errors"
	bsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

// AbiquoConfig ...
type AbiquoConfig struct {
	DatacenterRepository string
	Endpoint             string
	Username             string
	Password             string
}

// Config ...
type Config struct {
	Abiquo AbiquoConfig
	Agent  apiv1.AgentOptions
}

func newConfigFromPath(path string, fs bsys.FileSystem) (Config, error) {
	var config Config

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return config, berr.WrapErrorf(err, "Reading config '%s'", path)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, berr.WrapError(err, "Unmarshalling config")
	}

	err = config.Validate()
	if err != nil {
		return config, berr.WrapError(err, "Validating config")
	}

	return config, nil
}

// Validate ...
func (c Config) Validate() error {
	err := c.Abiquo.Validate()
	if err != nil {
		return berr.WrapError(err, "Validating Abiquo configuration")
	}

	err = c.Agent.Validate()
	if err != nil {
		return berr.WrapError(err, "Validating Agent configuration")
	}

	return nil
}

// Validate ...
func (c AbiquoConfig) Validate() error {
	if c.Endpoint == "" {
		return berr.Error("Must provide non-empty Endpoint")
	}

	if c.Username == "" {
		return berr.Error("Must provide non-empty Username")
	}

	if c.Password == "" {
		return berr.Error("Must provide non-empty Password")
	}

	if c.DatacenterRepository == "" {
		return berr.Error("Must provide non-empty DatacenterRepository")
	}

	return nil
}
