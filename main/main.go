package main

import (
	"flag"
	"os"

	"github.com/cloudfoundry/bosh-utils/fileutil"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/bosh-utils/system"
	bsys "github.com/cloudfoundry/bosh-utils/system"
	"github.com/cloudfoundry/bosh-utils/uuid"

	"github.com/agoca80/bosh-abiquo-cpi/actions"
	"github.com/cppforlife/bosh-cpi-go/rpc"
)

var configPathOpt = flag.String("configPath", "", "Path to configuration file")

func main() {
	logger, fs, cmdRunner, uuidGen := basicDeps()
	defer logger.HandlePanic("Main")

	flag.Parse()
	config, err := newConfigFromPath(*configPathOpt, fs)
	if err != nil {
		logger.Error("main", "Loading config %s", err.Error())
		os.Exit(1)
	}

	compressor := fileutil.NewTarballCompressor(cmdRunner, fs)

	cpiFactory := actions.NewFactory(
		fs,
		cmdRunner,
		uuidGen,
		compressor,
		logger,
		actions.Options{
			DatacenterRepository: config.Abiquo.DatacenterRepository,
			Password:             config.Abiquo.Password,
			Username:             config.Abiquo.Username,
			Endpoint:             config.Abiquo.Endpoint,
			Agent:                config.Agent,
		},
	)

	cli := rpc.NewFactory(logger).NewCLI(cpiFactory)

	err = cli.ServeOnce()
	if err != nil {
		logger.Error("main", "Serving once: %s", err)
		os.Exit(1)
	}
}

func basicDeps() (logger.Logger, bsys.FileSystem, bsys.CmdRunner, uuid.Generator) {
	var (
		uuidGen = uuid.NewGenerator()
		logger  = logger.NewWriterLogger(logger.LevelDebug, os.Stderr)
		runner  = system.NewExecCmdRunner(logger)
		fs      = system.NewOsFileSystem(logger)
	)

	return logger, fs, runner, uuidGen
}
